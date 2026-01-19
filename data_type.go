package edgeexpr

import (
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
)

// generate datatype enumeration
type DataType string

const (
	DataTypeBool    DataType = "Bool"
	DataTypeByte    DataType = "Byte"
	DataTypeWord    DataType = "Word"
	DataTypeDWord   DataType = "DWord"
	DataTypeInt8    DataType = "Int8"
	DataTypeUInt8   DataType = "UInt8"
	DataTypeInt16   DataType = "Int16"
	DataTypeUInt16  DataType = "UInt16"
	DataTypeInt32   DataType = "Int32"
	DataTypeUInt32  DataType = "UInt32"
	DataTypeInt64   DataType = "Int64"
	DataTypeUInt64  DataType = "UInt64"
	DataTypeFloat32 DataType = "Float32"
	DataTypeFloat64 DataType = "Float64"
	DataTypeString  DataType = "String"
)

func (dt DataType) String() string {
	return string(dt)
}

// DataTypeValidator is a validator for the "dataType" field enum values. It is called by the builders before save.
func DataTypeValidator(dt DataType) error {
	switch dt {
	case DataTypeBool, DataTypeByte, DataTypeWord, DataTypeDWord, DataTypeInt8, DataTypeUInt8, DataTypeInt16, DataTypeUInt16, DataTypeInt32, DataTypeUInt32, DataTypeInt64, DataTypeUInt64, DataTypeFloat32, DataTypeFloat64, DataTypeString:
		return nil
	default:
		return fmt.Errorf("data: invalid enum value for dataType field: %q", dt)
	}
}

func (DataType) Values() []string {
	return []string{
		string(DataTypeBool),
		string(DataTypeByte),
		string(DataTypeWord),
		string(DataTypeDWord),
		string(DataTypeInt8),
		string(DataTypeUInt8),
		string(DataTypeInt16),
		string(DataTypeUInt16),
		string(DataTypeInt32),
		string(DataTypeUInt32),
		string(DataTypeInt64),
		string(DataTypeUInt64),
		string(DataTypeFloat32),
		string(DataTypeFloat64),
		string(DataTypeString),
	}
}

func (dt DataType) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(string(dt)))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (dt *DataType) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", val)
	}
	*dt = DataType(str)
	if err := DataTypeValidator(*dt); err != nil {
		return fmt.Errorf("%s is not a valid DataType", str)
	}
	return nil
}

func ParseDataType(dt string) (DataType, int, error) {
	switch dt {
	case string(DataTypeBool):
		return DataTypeBool, 1, nil
	case string(DataTypeByte):
		return DataTypeByte, 1, nil
	case string(DataTypeWord):
		return DataTypeWord, 2, nil
	case string(DataTypeDWord):
		return DataTypeDWord, 4, nil
	case string(DataTypeInt8):
		return DataTypeInt8, 1, nil
	case string(DataTypeUInt8):
		return DataTypeUInt8, 1, nil
	case string(DataTypeInt16):
		return DataTypeInt16, 2, nil
	case string(DataTypeUInt16):
		return DataTypeUInt16, 2, nil
	case string(DataTypeInt32):
		return DataTypeInt32, 4, nil
	case string(DataTypeUInt32):
		return DataTypeUInt32, 4, nil
	case string(DataTypeInt64):
		return DataTypeInt64, 8, nil
	case string(DataTypeUInt64):
		return DataTypeUInt64, 8, nil
	case string(DataTypeFloat32):
		return DataTypeFloat32, 4, nil
	case string(DataTypeFloat64):
		return DataTypeFloat64, 8, nil
	case string(DataTypeString):
		return DataTypeString, 0, nil // String has no fixed size
	case "S5Time": //ms
		return DataTypeInt16, 2, nil
	case "Time": //ms
		return DataTypeInt32, 4, nil
	case "LTime": //ns
		return DataTypeInt64, 8, nil
	case "DTL":
		return DataTypeString, 12, nil
	case "Date":
		return DataTypeString, 2, nil
	case "Date_And_Time":
		return DataTypeString, 8, nil
	case "LDT":
		return DataTypeString, 8, nil
	case "LTime_Of_Day":
		return DataTypeString, 8, nil
	case "Time_Of_Day":
		return DataTypeString, 4, nil
	default:
		// for siemens like "WString[10]", "String[20]", etc.
		reg, _ := regexp.Compile(`^(W)?String\[(\d+)\]$`)
		match := reg.FindStringSubmatch(dt)
		if match != nil {
			ll, _ := strconv.Atoi(match[2])
			if match[1] == "W" {
				return DataTypeString, ll*2 + 4, nil
			}
			return DataTypeString, ll + 2, nil
		}
	}
	return "", 0, fmt.Errorf("unknown data type: %s", dt)
}

func ConvertToFloat64(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}

func (dt DataType) ConvertFromAny(value any) (any, error) {
	switch dt {
	case DataTypeBool:
		switch v := value.(type) {
		case bool:
			return v, nil
		case int, uint, uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64:
			return v != 0, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to bool", value)
		}
	case DataTypeInt8:
		switch v := value.(type) {
		case uint:
			if v > math.MaxInt8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int8: out of range", v, value)
			}
			return int8(v), nil
		case uint8:
			if v > math.MaxInt8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int8: out of range", v, value)
			}
			return int8(v), nil
		case uint16:
			if v > math.MaxInt8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int8: out of range", v, value)
			}
			return int8(v), nil
		case uint32:
			if v > math.MaxInt8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int8: out of range", v, value)
			}
			return int8(v), nil
		case uint64:
			if v > math.MaxInt8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int8: out of range", v, value)
			}
			return int8(v), nil
		case int:
			return int8(v), nil
		case int8:
			return v, nil
		case int16:
			return int8(v), nil
		case int32:
			return int8(v), nil
		case int64:
			return int8(v), nil
		case float32:
			if v < math.MinInt8 || v > math.MaxInt8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int8: out of range", v, value)
			}
			return int8(v), nil
		case float64:
			if v < math.MinInt8 || v > math.MaxInt8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int8: out of range", v, value)
			}
			return int8(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to int8", value)
		}
	case DataTypeInt16:
		switch v := value.(type) {
		case uint:
			return int16(v), nil
		case uint8:
			return int16(v), nil
		case uint16:
			return int16(v), nil
		case uint32:
			return int16(v), nil
		case uint64:
			return int16(v), nil
		case int:
			return int16(v), nil
		case int8:
			return int16(v), nil
		case int16:
			return v, nil
		case int32:
			return int16(v), nil
		case int64:
			return int16(v), nil
		case float32:
			if v < math.MinInt16 || v > math.MaxInt16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int16: out of range", v, value)
			}
			return int16(v), nil
		case float64:
			if v < math.MinInt16 || v > math.MaxInt16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int16: out of range", v, value)
			}
			return int16(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to int16", value)
		}
	case DataTypeInt32:
		switch v := value.(type) {
		case uint:
			return int32(v), nil
		case uint8:
			return int32(v), nil
		case uint16:
			return int32(v), nil
		case uint32:
			return int32(v), nil
		case uint64:
			return int32(v), nil
		case int:
			return int32(v), nil
		case int8:
			return int32(v), nil
		case int16:
			return int32(v), nil
		case int32:
			return v, nil
		case int64:
			return int32(v), nil
		case float32:
			if v < math.MinInt32 || v > math.MaxInt32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int32: out of range", v, value)
			}
			return int32(v), nil
		case float64:
			if v < math.MinInt32 || v > math.MaxInt32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int32: out of range", v, value)
			}
			return int32(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to int32", value)
		}
	case DataTypeInt64:
		switch v := value.(type) {
		case uint:
			return int64(v), nil
		case uint8:
			return int64(v), nil
		case uint16:
			return int64(v), nil
		case uint32:
			return int64(v), nil
		case uint64:
			return int64(v), nil
		case int:
			return int64(v), nil
		case int8:
			return int64(v), nil
		case int16:
			return int64(v), nil
		case int32:
			return int64(v), nil
		case int64:
			return v, nil
		case float32:
			if v < math.MinInt64 || v > math.MaxInt64 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int64: out of range", v, value)
			}
			return int64(v), nil
		case float64:
			if v < math.MinInt64 || v > math.MaxInt64 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to int64: out of range", v, value)
			}
			return int64(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to int64", value)
		}
	case DataTypeUInt8:
		switch v := value.(type) {
		case int:
			if v < 0 || v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		case int8:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		case int16:
			if v < 0 || v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		case int32:
			if v < 0 || v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		case int64:
			if v < 0 || v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		case uint:
			if v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		case uint8:
			return v, nil
		case uint16:
			if v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		case uint32:
			if v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		case uint64:
			if v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		case float32:
			if v < 0 || v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		case float64:
			if v < 0 || v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint8: out of range", v, value)
			}
			return uint8(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to uint8", value)
		}
	case DataTypeUInt16:
		switch v := value.(type) {
		case int:
			if v < 0 || v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint16: out of range", v, value)
			}
			return uint16(v), nil
		case int8:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint16: out of range", v, value)
			}
			return uint16(v), nil
		case int16:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint16: out of range", v, value)
			}
			return uint16(v), nil
		case int32:
			if v < 0 || v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint16: out of range", v, value)
			}
			return uint16(v), nil
		case int64:
			if v < 0 || v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint16: out of range", v, value)
			}
			return uint16(v), nil
		case uint:
			if v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint16: out of range", v, value)
			}
			return uint16(v), nil
		case uint8:
			return uint16(v), nil
		case uint16:
			return v, nil
		case uint32:
			if v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint16: out of range", v, value)
			}
			return uint16(v), nil
		case uint64:
			if v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint16: out of range", v, value)
			}
			return uint16(v), nil
		case float32:
			if v < 0 || v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint16: out of range", v, value)
			}
			return uint16(v), nil
		case float64:
			if v < 0 || v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint16: out of range", v, value)
			}
			return uint16(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to uint16", value)
		}
	case DataTypeUInt32:
		switch v := value.(type) {
		case int:
			if v < 0 || uint64(v) > math.MaxUint32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint32: out of range", v, value)
			}
			return uint32(v), nil
		case int8:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint32: out of range", v, value)
			}
			return uint32(v), nil
		case int16:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint32: out of range", v, value)
			}
			return uint32(v), nil
		case int32:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint32: out of range", v, value)
			}
			return uint32(v), nil
		case int64:
			if v < 0 || uint64(v) > math.MaxUint32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint32: out of range", v, value)
			}
			return uint32(v), nil
		case uint:
			if uint64(v) > math.MaxUint32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint32: out of range", v, value)
			}
			return uint32(v), nil
		case uint8:
			return uint32(v), nil
		case uint16:
			return uint32(v), nil
		case uint32:
			return v, nil
		case uint64:
			if v > math.MaxUint32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint32: out of range", v, value)
			}
			return uint32(v), nil
		case float32:
			if v < 0 || float64(v) > math.MaxUint32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint32: out of range", v, value)
			}
			return uint32(v), nil
		case float64:
			if v < 0 || v > math.MaxUint32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint32: out of range", v, value)
			}
			return uint32(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to uint32", value)
		}
	case DataTypeUInt64:
		switch v := value.(type) {
		case int:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint64: out of range", v, value)
			}
			return uint64(v), nil
		case int8:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint64: out of range", v, value)
			}
			return uint64(v), nil
		case int16:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint64: out of range", v, value)
			}
			return uint64(v), nil
		case int32:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint64: out of range", v, value)
			}
			return uint64(v), nil
		case int64:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint64: out of range", v, value)
			}
			return uint64(v), nil
		case uint:
			return uint64(v), nil
		case uint8:
			return uint64(v), nil
		case uint16:
			return uint64(v), nil
		case uint32:
			return uint64(v), nil
		case uint64:
			return v, nil
		case float32:
			if v < 0 || float64(v) > math.MaxUint64 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint64: out of range", v, value)
			}
			return uint64(v), nil
		case float64:
			if v < 0 || v > math.MaxUint64 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to uint64: out of range", v, value)
			}
			return uint64(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to uint64", value)
		}
	case DataTypeFloat32:
		switch v := value.(type) {
		case float32:
			return v, nil
		case float64:
			if v > math.MaxFloat32 || v < -math.MaxFloat32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to float32: out of range", v, value)
			}
			return float32(v), nil
		case int:
			if float64(v) > math.MaxFloat32 || float64(v) < -math.MaxFloat32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to float32: out of range", v, value)
			}
			return float32(v), nil
		case int8:
			return float32(v), nil
		case int16:
			return float32(v), nil
		case int32:
			return float32(v), nil
		case int64:
			if float64(v) > math.MaxFloat32 || float64(v) < -math.MaxFloat32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to float32: out of range", v, value)
			}
			return float32(v), nil
		case uint:
			if float64(v) > math.MaxFloat32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to float32: out of range", v, value)
			}
			return float32(v), nil
		case uint8:
			return float32(v), nil
		case uint16:
			return float32(v), nil
		case uint32:
			return float32(v), nil
		case uint64:
			if float64(v) > math.MaxFloat32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to float32: out of range", v, value)
			}
			return float32(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to float32", value)
		}
	case DataTypeFloat64:
		switch v := value.(type) {
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case int:
			return float64(v), nil
		case int8:
			return float64(v), nil
		case int16:
			return float64(v), nil
		case int32:
			return float64(v), nil
		case int64:
			return float64(v), nil
		case uint:
			return float64(v), nil
		case uint8:
			return float64(v), nil
		case uint16:
			return float64(v), nil
		case uint32:
			return float64(v), nil
		case uint64:
			return float64(v), nil
		default:
			return nil, fmt.Errorf("cannot convert %T to float64", value)
		}
	case DataTypeString:
		switch v := value.(type) {
		case string:
			return v, nil
		case []byte:
			return string(v), nil
		default:
			return fmt.Sprintf("%v", value), nil
		}
	case DataTypeByte:
		switch v := value.(type) {
		case []byte:
			if len(v) > 1 {
				return nil, fmt.Errorf("cannot convert %T to [1]byte: too long", value)
			}
			var arr [1]byte
			copy(arr[:], v)
			return arr, nil
		case [1]byte:
			return v, nil
		case uint8:
			return [1]byte{v}, nil
		case uint16:
			if v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [1]byte: out of range", v, value)
			}
			return [1]byte{byte(v)}, nil
		case uint32:
			if v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [1]byte: out of range", v, value)
			}
			return [1]byte{byte(v)}, nil
		case uint64:
			if v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [1]byte: out of range", v, value)
			}
			return [1]byte{byte(v)}, nil
		case uint:
			if v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [1]byte: out of range", v, value)
			}
			return [1]byte{byte(v)}, nil
		case int8:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [1]byte: out of range", v, value)
			}
			return [1]byte{byte(v)}, nil
		case int16:
			if v < 0 || v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [1]byte: out of range", v, value)
			}
			return [1]byte{byte(v)}, nil
		case int32:
			if v < 0 || v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [1]byte: out of range", v, value)
			}
			return [1]byte{byte(v)}, nil
		case int64:
			if v < 0 || v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [1]byte: out of range", v, value)
			}
			return [1]byte{byte(v)}, nil
		case int:
			if v < 0 || v > math.MaxUint8 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [1]byte: out of range", v, value)
			}
			return [1]byte{byte(v)}, nil
		case string:
			if len(v) > 1 {
				return nil, fmt.Errorf("cannot convert %T to [1]byte: string too long", value)
			}
			var arr [1]byte
			copy(arr[:], v)
			return arr, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to [1]byte", value)
		}
	case DataTypeWord:
		switch v := value.(type) {
		case []byte:
			if len(v) > 2 {
				return nil, fmt.Errorf("cannot convert %T to [2]byte: too long", value)
			}
			var arr [2]byte
			copy(arr[:], v)
			return arr, nil
		case [2]byte:
			return v, nil
		case uint8:
			return [2]byte{v, 0}, nil
		case uint16:
			return [2]byte{byte(v), byte(v >> 8)}, nil
		case uint32:
			if v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [2]byte: out of range", v, value)
			}
			return [2]byte{byte(v), byte(v >> 8)}, nil
		case uint64:
			if v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [2]byte: out of range", v, value)
			}
			return [2]byte{byte(v), byte(v >> 8)}, nil
		case uint:
			if v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [2]byte: out of range", v, value)
			}
			return [2]byte{byte(v), byte(v >> 8)}, nil

		case int8:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [2]byte: out of range", v, value)
			}
			return [2]byte{byte(v), 0}, nil
		case int16:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [2]byte: out of range", v, value)
			}
			return [2]byte{byte(v), byte(v >> 8)}, nil
		case int32:
			if v < 0 || v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [2]byte: out of range", v, value)
			}
			return [2]byte{byte(v), byte(v >> 8)}, nil
		case int64:
			if v < 0 || v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [2]byte: out of range", v, value)
			}
			return [2]byte{byte(v), byte(v >> 8)}, nil
		case int:
			if v < 0 || v > math.MaxUint16 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [2]byte: out of range", v, value)
			}
			return [2]byte{byte(v), byte(v >> 8)}, nil
		case string:
			if len(v) > 2 {
				return nil, fmt.Errorf("cannot convert %T to [2]byte: string too long", value)
			}
			var arr [2]byte
			copy(arr[:], v)
			return arr, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to [2]byte", value)
		}
	case DataTypeDWord:
		switch v := value.(type) {
		case []byte:
			if len(v) > 4 {
				return nil, fmt.Errorf("cannot convert %T to [4]byte: too long", value)
			}
			var arr [4]byte
			copy(arr[:], v)
			return arr, nil
		case [4]byte:
			return v, nil
		case uint8:
			return [4]byte{v, 0, 0, 0}, nil
		case uint16:
			return [4]byte{byte(v), byte(v >> 8), 0, 0}, nil
		case uint32:
			return [4]byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}, nil
		case uint64:
			if v > math.MaxUint32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [4]byte: out of range", v, value)
			}
			return [4]byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}, nil
		case uint:
			if uint64(v) > math.MaxUint32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [4]byte: out of range", v, value)
			}
			return [4]byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}, nil
		case int8:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [4]byte: out of range", v, value)
			}
			return [4]byte{byte(v), 0, 0, 0}, nil
		case int16:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [4]byte: out of range", v, value)
			}
			return [4]byte{byte(v), byte(v >> 8), 0, 0}, nil
		case int32:
			if v < 0 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [4]byte: out of range", v, value)
			}
			return [4]byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}, nil
		case int64:
			if v < 0 || uint64(v) > math.MaxUint32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [4]byte: out of range", v, value)
			}
			return [4]byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}, nil
		case int:
			if v < 0 || uint64(v) > math.MaxUint32 {
				return nil, fmt.Errorf("cannot convert %v (type %T) to [4]byte: out of range", v, value)
			}
			return [4]byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}, nil
		case string:
			if len(v) > 4 {
				return nil, fmt.Errorf("cannot convert %T to [4]byte: string too long", value)
			}
			var arr [4]byte
			copy(arr[:], v)
			return arr, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to [4]byte", value)
		}
	default:
		return nil, fmt.Errorf("unsupported data type: %v", dt)
	}
}

func ConvertToBytes(value any) ([]byte, error) {
	switch v := value.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case [1]byte:
		return v[:], nil
	case [2]byte:
		return v[:], nil
	case [4]byte:
		return v[:], nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}
