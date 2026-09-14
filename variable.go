package edgeexpr

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/expr-lang/expr/vm"
)

type Variable struct {
	Key           string         `json:"key"`
	Connection    string         `json:"connection"`
	Address       string         `json:"address"`
	Script        string         `json:"script"`
	DiffThreshold *float64       `json:"diff_threshold,omitempty"` // Optional threshold for change detection, in the same unit as the variable
	PctThreshold  *float64       `json:"pct_threshold,omitempty"`  // Optional percentage threshold for change detection, in the same unit as the variable
	Scale         *float64       `json:"scale,omitempty"`          // Optional scale factor for the variable value
	Offset        *float64       `json:"offset,omitempty"`         // Optional offset for the variable value
	Writable      bool           `json:"writable,omitempty"`       // Optional flag to indicate if the variable is writable
	DataTypeStr   string         `json:"data_type"`
	DataType      DataType       `json:"-"`
	Bytes         int            `json:"-"` // Number of bytes for the data type, derived from DataType
	PublishCycle  *time.Duration `json:"-"`
	CacheDuration *time.Duration `json:"-"`

	Cache      *Cache      `json:"-"`
	LatestPush *Point      `json:"-"`
	Program    *vm.Program `json:"-"`
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
	}
	v.Cache = NewCache(v.CacheDuration) // Create cache instance based on DataType and CacheDuration
	return nil
}

func Check[T any](v *T) T {
	return *v
}

func (v *Variable) Hash() string {
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
	if v.CacheDuration != nil {
		hash.Write([]byte(v.CacheDuration.String()))
	}
	if v.PublishCycle != nil {
		hash.Write([]byte(v.PublishCycle.String()))
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (v *Variable) Read() (any, *time.Time) {
	if v.Cache == nil {
		return nil, nil
	}
	return v.Cache.Value(), v.Cache.Timestamp()
}

func (v *Variable) ValueUnScale(value interface{}) interface{} {
	switch val := value.(type) {
	case float64:
		if v.Scale != nil && *v.Scale != 0 {
			val /= *v.Scale
		}
		if v.Offset != nil {
			val -= *v.Offset
		}
		return val
	case float32:
		if v.Scale != nil && *v.Scale != 0 {
			val /= float32(*v.Scale)
		}
		if v.Offset != nil {
			val -= float32(*v.Offset)
		}
		return val
	default:
		return value
	}
}

func (v *Variable) WriteValue(value any, t *time.Time) error {
	v.Cache.AddPoint(value, t)
	return nil
}
