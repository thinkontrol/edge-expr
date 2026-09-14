package edgeexpr

import (
	"math"
	"reflect"
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
	publishCycle := int64(*v.PublishCycle)
	times := publishCycle / gcd
	changed := v.ChangedWithLatestPushValue()
	now := time.Now()
	if (publishCycle <= 0 && changed) || (times > 0 && i%times == 0) {
		if latestValue := v.Cache.Value(); latestValue != nil {
			pushValues = append(pushValues, &PushValue{
				Value:     latestValue,
				Timestamp: &now,
			})
			v.LatestPush = &Point{
				Value:     latestValue,
				Timestamp: &now,
			}
		}
	}
	return pushValues
}

func (v *Variable) ChangedWithLatestPushValue() bool {
	if v.Cache == nil {
		return false
	}
	if v.LatestPush == nil {
		return true
	}
	// changed := reflect.DeepEqual(v.Cache.Value(), v.LatestPush.Value)
	var changed bool
	lastPushTime := v.LatestPush.Timestamp

	switch cache := v.Cache.Value().(type) {
	case bool, string:
		changed = cache != v.LatestPush.Value
	case []byte:
		changed = reflect.DeepEqual(cache, v.LatestPush.Value)
	default:
		fv1, err := ConvertToFloat64(cache)
		fv2, err2 := ConvertToFloat64(v.LatestPush.Value)
		if err != nil || err2 != nil {
			changed = reflect.DeepEqual(cache, v.LatestPush.Value)
		} else {
			if v.DiffThreshold != nil {
				changed = math.Abs(fv1-fv2) >= *v.DiffThreshold
				break
			}
			if v.PctThreshold != nil {
				percentageChange := lo.Ternary(fv2 == 0, lo.Ternary(fv1 == 0, 0, math.MaxFloat64), ((fv1-fv2)/fv2)*100)
				changed = math.Abs(percentageChange) >= *v.PctThreshold
				break
			}
			changed = fv1 != fv2
		}
	}

	if changed {
		return true
	}
	if v.PublishCycle != nil && *v.PublishCycle < 0 {
		return lastPushTime == nil || time.Since(*lastPushTime) > (*v.PublishCycle).Abs()
	}
	return false
}
