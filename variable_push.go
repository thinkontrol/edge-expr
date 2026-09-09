package edgeexpr

import (
	"math"
	"time"

	"github.com/samber/lo"
)

func (v *Variable) GetPushValues(gcd, i int64) []*PushValue {
	var pushValues []*PushValue
	if v.PublishCycle == nil {
		return pushValues
	}
	if v.Cache == nil {
		return pushValues
	}
	// if !v.TimestampUpdated() {
	// 	return pushValues
	// }
	publishCycle := int64(*v.PublishCycle)
	times := publishCycle / gcd
	changed := v.ChangedWithLatestPushValue()
	if (publishCycle <= 0 && changed) || (times > 0 && i%times == 0) {
		switch cache := v.Cache.(type) {
		case *Cache[float64]:
			if pushValue := cache.PushValue(); pushValue != nil {
				if changed && len(cache.Points) >= 2 {
					if p, ok := v.LatestPush.(Point[float64]); ok {
						if p.Timestamp != nil && cache.Points[len(cache.Points)-2].Timestamp != nil && !p.Timestamp.Equal(*cache.Points[len(cache.Points)-2].Timestamp) {
							pushValues = append(pushValues, &PushValue{
								// Key:       v.Key,
								Value:     cache.Points[len(cache.Points)-2].Value,
								Timestamp: cache.Points[len(cache.Points)-2].Timestamp,
							})
						}
					}
				}
				pushValues = append(pushValues, pushValue)
				v.LatestPush = cache.Points[len(cache.Points)-1]
			}
		case *Cache[bool]:
			if pushValue := cache.PushValue(); pushValue != nil {
				pushValues = append(pushValues, pushValue)
				v.LatestPush = cache.Points[len(cache.Points)-1]
			}
		case *Cache[string]:
			if pushValue := cache.PushValue(); pushValue != nil {
				pushValues = append(pushValues, pushValue)
				v.LatestPush = cache.Points[len(cache.Points)-1]
			}
		case *Cache[[]byte]:
			if pushValue := cache.PushValue(); pushValue != nil {
				pushValues = append(pushValues, pushValue)
				v.LatestPush = cache.Points[len(cache.Points)-1]
			}
			// Supported cache types
		default:
			// Unsupported cache type
		}
	}
	return pushValues
}

func (v *Variable) TimestampUpdated() bool {
	if v.Cache == nil {
		return false
	}
	if v.LatestPush == nil {
		return true
	}
	switch cache := v.Cache.(type) {
	case *Cache[float64]:
		latestPush, ok := v.LatestPush.(Point[float64])
		if !ok {
			return true
		}
		return cache.Timestamp() != nil && (latestPush.Timestamp == nil || !cache.Timestamp().Equal(*latestPush.Timestamp))
	case *Cache[bool]:
		latestPush, ok := v.LatestPush.(Point[bool])
		if !ok {
			return true
		}
		return cache.Timestamp() != nil && (latestPush.Timestamp == nil || !cache.Timestamp().Equal(*latestPush.Timestamp))
	case *Cache[string]:
		latestPush, ok := v.LatestPush.(Point[string])
		if !ok {
			return true
		}
		return cache.Timestamp() != nil && (latestPush.Timestamp == nil || !cache.Timestamp().Equal(*latestPush.Timestamp))
	case *Cache[[]byte]:
		latestPush, ok := v.LatestPush.(Point[[]byte])
		if !ok {
			return true
		}
		return cache.Timestamp() != nil && (latestPush.Timestamp == nil || !cache.Timestamp().Equal(*latestPush.Timestamp))
		// Supported cache types
	default:
		return true // Unsupported cache type
	}
}

func (v *Variable) ChangedWithLatestPushValue() bool {
	if v.Cache == nil {
		return false
	}
	if v.LatestPush == nil {
		return true
	}
	var changed bool
	var lastPushTime *time.Time
	switch cache := v.Cache.(type) {
	case *Cache[float64]:
		latestPush, ok := v.LatestPush.(Point[float64])
		if !ok {
			changed = true
			break
		}
		lastPushTime = latestPush.Timestamp
		if v.DiffThreshold != nil {
			changed = math.Abs(cache.Value()-latestPush.Value) >= *v.DiffThreshold
			break
		}
		if v.PctThreshold != nil {
			percentageChange := lo.Ternary(latestPush.Value == 0, lo.Ternary(cache.Value() == 0, 0, math.MaxFloat64), ((cache.Value()-latestPush.Value)/latestPush.Value)*100)
			changed = math.Abs(percentageChange) >= *v.PctThreshold
			break
		}
		changed = cache.Value() != latestPush.Value
	case *Cache[bool]:
		latestPush, ok := v.LatestPush.(Point[bool])
		if !ok {
			changed = true
		} else {
			lastPushTime = latestPush.Timestamp
			changed = cache.Value() != latestPush.Value
		}
	case *Cache[string]:
		latestPush, ok := v.LatestPush.(Point[string])
		if !ok {
			changed = true
		} else {
			lastPushTime = latestPush.Timestamp
			changed = cache.Value() != latestPush.Value
		}
	case *Cache[[]byte]:
		latestPush, ok := v.LatestPush.(Point[[]byte])
		if !ok {
			changed = true
			break
		}
		lastPushTime = latestPush.Timestamp
		if len(cache.Value()) != len(latestPush.Value) {
			changed = true
		} else {
			for i := range cache.Value() {
				if cache.Value()[i] != latestPush.Value[i] {
					changed = true
					break
				}
			}
		}
		// Supported cache types
	default:
		return true // Unsupported cache type
	}
	if changed {
		return true
	}
	if v.PublishCycle != nil && *v.PublishCycle < 0 {
		return lastPushTime == nil || time.Since(*lastPushTime) > (*v.PublishCycle).Abs()
	}
	return false
}
