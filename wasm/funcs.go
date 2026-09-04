package main

import (
	"encoding/base64"
	"fmt"
	"math"

	"github.com/expr-lang/expr"
)

var BitFunc = expr.Function(
	"bitOn",
	func(params ...any) (any, error) {
		if len(params) != 2 {
			return nil, fmt.Errorf("bitOn() 需要 2 个参数")
		}

		var bb []byte

		switch v := params[0].(type) {
		case string:
			// 尝试 Base64 解码，如果失败则回退到原始 string 字节流
			decoded, err := base64.StdEncoding.DecodeString(v)
			if err == nil {
				bb = decoded
			} else {
				bb = []byte(v)
			}
		case []byte:
			bb = v
		default:
			return nil, fmt.Errorf("不支持的类型: %T", params[0])
		}

		pos, err := bitPosition(params[1])
		if err != nil {
			return nil, err
		}

		if pos < 0 || pos/8 >= int64(len(bb)) {
			return nil, fmt.Errorf("position %d 超出范围 (字节长度: %d)", pos, len(bb))
		}

		// 取出对应的 Bit
		return (bb[pos/8] & (1 << (pos % 8))) != 0, nil
	},
)

func bitPosition(value any) (int64, error) {
	switch position := value.(type) {
	case int:
		return int64(position), nil
	case int8:
		return int64(position), nil
	case int16:
		return int64(position), nil
	case int32:
		return int64(position), nil
	case int64:
		return position, nil
	case uint:
		if uint64(position) <= math.MaxInt64 {
			return int64(position), nil
		}
	case uint8:
		return int64(position), nil
	case uint16:
		return int64(position), nil
	case uint32:
		return int64(position), nil
	case uint64:
		if position <= math.MaxInt64 {
			return int64(position), nil
		}
	}

	return 0, fmt.Errorf("position 必须为整数")
}
