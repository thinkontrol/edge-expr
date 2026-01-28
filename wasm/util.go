//go:build js && wasm
// +build js,wasm

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/expr-lang/expr"
	edgeexpr "github.com/thinkontrol/edge-expr"
)

func ValidScript(plcVar map[string]string, script string) (string, error) {
	env := make(map[string]any)
	for key, dtStr := range plcVar {
		parts := strings.Split(key, ".")
		if len(parts) == 0 {
			return "", fmt.Errorf("invalid plc var key: %s", key)
		}
		cur := env
		for i, k := range parts {
			if i == len(parts)-1 {
				dt, _, err := edgeexpr.ParseDataType(dtStr)
				if err != nil {
					return "", err
				}
				rv, err := dt.GenerateRandomValue()
				if err != nil {
					return "", err
				}
				cur[k] = rv
			} else {
				if next, ok := cur[k].(map[string]any); ok {
					cur = next
				} else {
					newm := make(map[string]any)
					cur[k] = newm
					cur = newm
				}
			}
		}

	}

	js.Global().Get("console").Call("log", marshalJSON(env))
	js.Global().Get("console").Call("log", script)

	program, err := expr.Compile(script, expr.Env(env))
	if err != nil {
		return "", err
	}
	out, err := expr.Run(program, env)
	if err != nil {
		return "", err
	}
	dt := InferTypeName(out)
	if dt == "" {
		return "", fmt.Errorf("unsupported result data type")
	}
	return dt, nil
}

func unmarshalJSON(src js.Value, dst interface{}) error {
	str := js.Global().Get("JSON").Call("stringify", src).String()
	return json.Unmarshal([]byte(str), &dst)
}

func marshalJSON(src interface{}) js.Value {
	if src == nil {
		return js.Null()
	}
	b, e := json.Marshal(src)
	if e != nil {
		return js.Null()
	}
	return js.Global().Get("JSON").Call("parse", string(b))
}

func parseNumber(s string) (uint64, error) {
	if len(s) > 2 && (strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X")) {
		return strconv.ParseUint(s[2:], 16, 16)
	}
	return strconv.ParseUint(s, 10, 16)
}

func InferTypeName(value any) string {
	switch value.(type) {
	case bool:
		return "bool"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "int"
	case float32, float64:
		return "float"
	case string:
		return "string"
	default:
		return ""
	}
}
