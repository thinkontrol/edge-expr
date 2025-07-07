package edgeexpr

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/expr-lang/expr/vm"
)

type Variable struct {
	Key           string   `json:"key"`
	Connection    string   `json:"connection"`
	Address       string   `json:"address"`
	Script        string   `json:"script"`
	DiffThreshold *float64 `json:"diff_threshold,omitempty"` // Optional threshold for change detection, in the same unit as the variable
	PctThreshold  *float64 `json:"pct_threshold,omitempty"`  // Optional percentage threshold for change detection, in the same unit as the variable
	Scale         *float64 `json:"scale,omitempty"`          // Optional scale factor for the variable value
	Offset        *float64 `json:"offset,omitempty"`         // Optional offset for the variable value
	Writable      bool     `json:"writable,omitempty"`       // Optional flag to indicate if the variable is writable
	AsTag         bool     `json:"as_tag,omitempty"`         // Optional flag to indicate if the variable should be treated as a tag
	AsEvent       bool     `json:"as_event,omitempty"`       // Optional flag to indicate if the variable should be treated as an event
	DataTypeStr   string   `json:"data_type"`
	DataType      DataType
	Bytes         int
	PublishCycle  *time.Duration
	CacheDuration *time.Duration // Store cache duration instead of cache instance

	Cache   any
	Program *vm.Program
	// Cache instances can be created externally when needed
	// This allows the Variable to be non-generic while still supporting caching
}

func (v *Variable) MarshalJSON() ([]byte, error) {
	type Alias Variable

	// Prepare the auxiliary struct with string representations
	aux := &struct {
		*Alias
		PublishCycleStr  string `json:"publish_cycle,omitempty"`
		CacheDurationStr string `json:"cache_duration,omitempty"`
	}{
		Alias: (*Alias)(v),
	}

	// Convert PublishCycle to string if it exists
	if v.PublishCycle != nil {
		aux.PublishCycleStr = v.PublishCycle.String()
	}

	// Convert cache duration to string if cache exists
	if v.CacheDuration != nil {
		aux.CacheDurationStr = v.CacheDuration.String()
	}

	return json.Marshal(aux)
}

func (v *Variable) UnmarshalJSON(data []byte) error {
	type Alias Variable
	aux := &struct {
		*Alias
		PublishCycleStr  string `json:"publish_cycle"`
		CacheDurationStr string `json:"cache_duration"`
	}{
		Alias: (*Alias)(v),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	var err error
	v.DataType, v.Bytes, err = ParseDataType(aux.DataTypeStr)
	if v.Connection != "" && err != nil {
		return err
	}

	if v.AsTag && v.DataType != DataTypeString {
		return fmt.Errorf("variable %s with data type %s cannot be used as a tag", v.Key, v.DataType)
	}
	if v.AsEvent && v.DataType != DataTypeBool {
		return fmt.Errorf("variable %s with data type %s cannot be used as an event", v.Key, v.DataType)
	}

	// Parse PublishCycle to time.Duration and set publishCycle
	if aux.PublishCycleStr != "" {
		if duration, err := time.ParseDuration(aux.PublishCycleStr); err == nil {
			v.PublishCycle = &duration
		} else {
			return fmt.Errorf("invalid publish_cycle format: %v", err)
		}
	}
	// Parse Cache to time.Duration and set cache duration
	if aux.CacheDurationStr != "" {
		if duration, err := time.ParseDuration(aux.CacheDurationStr); err == nil {
			v.CacheDuration = &duration
		} else {
			return fmt.Errorf("invalid cache format: %v", err)
		}
	} else {
		defaultDuration := time.Minute
		v.CacheDuration = &defaultDuration // Default cache duration if not specified
	}
	v.Cache = v.createCache() // Create cache instance based on DataType and CacheDuration
	return nil
}

func (v *Variable) Hash() string {
	// Implement a hash function to generate a unique identifier for the variable
	hash := md5.New()
	hash.Write([]byte(v.Key))
	hash.Write([]byte(v.Connection))
	hash.Write([]byte(v.Address))
	hash.Write([]byte(v.Script))
	hash.Write([]byte(v.DataTypeStr))
	if v.DiffThreshold != nil {
		hash.Write([]byte(fmt.Sprintf("%0.8f", *v.DiffThreshold)))
	}
	if v.PctThreshold != nil {
		hash.Write([]byte(fmt.Sprintf("%0.8f", *v.PctThreshold)))
	}
	if v.Scale != nil {
		hash.Write([]byte(fmt.Sprintf("%0.8f", *v.Scale)))
	}
	if v.Offset != nil {
		hash.Write([]byte(fmt.Sprintf("%0.8f", *v.Offset)))
	}
	hash.Write([]byte(fmt.Sprintf("%t", v.Writable)))
	hash.Write([]byte(fmt.Sprintf("%t", v.AsTag)))
	hash.Write([]byte(fmt.Sprintf("%t", v.AsEvent)))
	if v.CacheDuration != nil {
		hash.Write([]byte(v.CacheDuration.String()))
	}
	if v.PublishCycle != nil {
		hash.Write([]byte(v.PublishCycle.String()))
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (v *Variable) Read() (any, bool, *time.Time) {
	var changed bool
	switch v.DataType {
	case DataTypeFloat32, DataTypeFloat64, DataTypeInt8, DataTypeUInt8, DataTypeInt16, DataTypeUInt16,
		DataTypeInt32, DataTypeUInt32, DataTypeInt64, DataTypeUInt64:
		if cache, ok := v.Cache.(*Cache[float64]); ok {
			if v.DiffThreshold != nil {
				changed, _ = cache.DiffExceeds(*v.DiffThreshold)
			}
			if v.PctThreshold != nil {
				changed, _ = cache.PctChangeExceeds(*v.PctThreshold)
			}
			return cache.Value(), changed, cache.Timestamp()
		}
	case DataTypeBool:
		if cache, ok := v.Cache.(*Cache[bool]); ok {
			return cache.Value(), cache.Changed(), cache.Timestamp()
		}
	case DataTypeString:
		if cache, ok := v.Cache.(*Cache[string]); ok {
			return cache.Value(), cache.Changed(), cache.Timestamp()
		}
	case DataTypeByte, DataTypeWord, DataTypeDWord:
		if cache, ok := v.Cache.(*Cache[[]byte]); ok {
			return cache.Value(), cache.Changed(), cache.Timestamp()
		}
	default:
		return nil, false, nil
	}
	return nil, false, nil // Unsupported data type or cache type mismatch
}

// func (v *Variable) Changed() bool {
// 	switch v.DataType {
// 	case DataTypeFloat32, DataTypeFloat64, DataTypeInt8, DataTypeUInt8, DataTypeInt16, DataTypeUInt16,
// 		DataTypeInt32, DataTypeUInt32, DataTypeInt64, DataTypeUInt64:
// 		if cache, ok := v.Cache.(*Cache[float64]); ok {
// 			if v.DiffThreshold != nil {
// 				v, _ := cache.DiffExceeds(*v.DiffThreshold)
// 				return v
// 			}
// 			if v.PctThreshold != nil {
// 				v, _ := cache.PctChangeExceeds(*v.PctThreshold)
// 				return v
// 			}
// 		}
// 	case DataTypeBool:
// 		if cache, ok := v.Cache.(*Cache[bool]); ok {
// 			return cache.Changed()
// 		}
// 	case DataTypeString:
// 		if cache, ok := v.Cache.(*Cache[string]); ok {
// 			return cache.Changed()
// 		}
// 	case DataTypeByte, DataTypeChar, DataTypeWord, DataTypeDWord:
// 		if cache, ok := v.Cache.(*Cache[[]byte]); ok {
// 			return cache.Changed()
// 		}
// 	default:
// 		return false
// 	}
// 	return false
// }

func (v *Variable) WriteValue(value any, t *time.Time) error {
	switch v.DataType {
	case DataTypeFloat32, DataTypeFloat64, DataTypeInt8, DataTypeUInt8, DataTypeInt16, DataTypeUInt16,
		DataTypeInt32, DataTypeUInt32, DataTypeInt64, DataTypeUInt64:
		floatValue, err := ConvertToFloat64(value)
		if err != nil {
			return err
		}
		if v.Scale != nil {
			floatValue *= *v.Scale
		}
		if v.Offset != nil {
			floatValue += *v.Offset
		}
		cache, ok := v.Cache.(*Cache[float64])
		if !ok {
			return fmt.Errorf("cache type mismatch for variable %s, expected Cache[float64]", v.Key)
		}
		cache.AddPoint(floatValue, t)
	case DataTypeBool:
		boolValue, err := v.DataType.ConvertFromAny(value)
		if err != nil {
			return fmt.Errorf("failed to convert value to bool for variable %s: %v", v.Key, err)
		}
		cache, ok := v.Cache.(*Cache[bool])
		if !ok {
			return fmt.Errorf("cache type mismatch for variable %s, expected Cache[bool]", v.Key)
		}
		cache.AddPoint(boolValue.(bool), t)
	case DataTypeString:
		stringValue, err := v.DataType.ConvertFromAny(value)
		if err != nil {
			return fmt.Errorf("failed to convert value to string for variable %s: %v", v.Key, err)
		}
		cache, ok := v.Cache.(*Cache[string])
		if !ok {
			return fmt.Errorf("cache type mismatch for variable %s, expected Cache[string]", v.Key)
		}
		cache.AddPoint(stringValue.(string), t)
	case DataTypeByte, DataTypeWord, DataTypeDWord:
		_bytesValue, err := v.DataType.ConvertFromAny(value)
		if err != nil {
			return fmt.Errorf("failed to convert value to bytes for variable %s: %v", v.Key, err)
		}
		bytesValue, err := ConvertToBytes(_bytesValue)
		if err != nil {
			return fmt.Errorf("failed to convert value to bytes for variable %s: %v", v.Key, err)
		}
		cache, ok := v.Cache.(*Cache[[]byte])
		if !ok {
			return fmt.Errorf("cache type mismatch for variable %s, expected Cache[[]byte]", v.Key)
		}
		cache.AddPoint(bytesValue, t)
	default:
		return fmt.Errorf("unsupported data type %s for writing value", v.DataType)
	}
	return nil
}

func (v *Variable) createCache() any {
	// 根据 DataType 创建相应类型的缓存
	switch v.DataType {
	case DataTypeFloat32, DataTypeFloat64, DataTypeInt8, DataTypeUInt8, DataTypeInt16, DataTypeUInt16,
		DataTypeInt32, DataTypeUInt32, DataTypeInt64, DataTypeUInt64:
		if v.CacheDuration != nil {
			return NewCache[float64](*v.CacheDuration)
		}
		return NewCache[float64](time.Minute)
	case DataTypeBool:
		if v.CacheDuration != nil {
			return NewCache[bool](*v.CacheDuration)
		}
		return NewCache[bool](time.Minute)
	case DataTypeString:
		if v.CacheDuration != nil {
			return NewCache[string](*v.CacheDuration)
		}
		return NewCache[string](time.Minute)
	case DataTypeByte, DataTypeWord, DataTypeDWord:
		if v.CacheDuration != nil {
			return NewCache[[]byte](*v.CacheDuration)
		}
		return NewCache[[]byte](time.Minute)
	default:
		return nil // Unsupported data type for caching
	}
}
