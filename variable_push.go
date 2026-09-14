package edgeexpr

import (
	"reflect"
	"time"
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
	changed := reflect.DeepEqual(v.Cache.Value(), v.LatestPush.Value)
	lastPushTime := v.LatestPush.Timestamp
	if changed {
		return true
	}
	if v.PublishCycle != nil && *v.PublishCycle < 0 {
		return lastPushTime == nil || time.Since(*lastPushTime) > (*v.PublishCycle).Abs()
	}
	return false
}
