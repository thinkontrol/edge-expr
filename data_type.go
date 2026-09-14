package edgeexpr

import (
	"fmt"
	"io"
	"math"
	"math/rand"
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

func (dt DataType) GenerateRandomValue() (any, error) {
	switch dt {
	case DataTypeBool:
		return rand.Intn(2) == 1, nil
	case DataTypeInt8:
		return int8(rand.Intn(math.MaxInt8 + 1)), nil
	case DataTypeInt16:
		return int16(rand.Intn(math.MaxInt16 + 1)), nil
	case DataTypeInt32:
		return int32(rand.Int31()), nil
	case DataTypeInt64:
		return rand.Int63(), nil
	case DataTypeUInt8:
		return uint8(rand.Intn(math.MaxUint8 + 1)), nil
	case DataTypeUInt16:
		return uint16(rand.Intn(math.MaxUint16 + 1)), nil
	case DataTypeUInt32:
		return uint32(rand.Uint32()), nil
	case DataTypeUInt64:
		return rand.Uint64(), nil
	case DataTypeFloat32:
		return rand.Float32() * math.MaxFloat32, nil
	case DataTypeFloat64:
		return rand.Float64() * math.MaxFloat64, nil
	case DataTypeString:
		length := rand.Intn(20) + 1 // Random length between 1 and 20
		bytes := make([]byte, length)
		for i := 0; i < length; i++ {
			bytes[i] = byte(rand.Intn(26) + 97) // a-z
		}
		return string(bytes), nil
	case DataTypeByte:
		var arr [1]byte
		rand.Read(arr[:])
		return arr[:], nil
	case DataTypeWord:
		var arr [2]byte
		rand.Read(arr[:])
		return arr[:], nil
	case DataTypeDWord:
		var arr [4]byte
		rand.Read(arr[:])
		return arr[:], nil
	default:
		return nil, fmt.Errorf("unsupported data type for random generation: %v", dt)
	}
}
