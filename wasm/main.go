//go:build js && wasm
// +build js,wasm

package main

import (
	_ "crypto/sha512"
	"syscall/js"
)

func main() {

	done := make(chan struct{})
	js.Global().Set("wasmValidScript", js.FuncOf(wasmValidScript))
	<-done
}

type ValidationResult struct {
	Error    string `json:"error"`
	ResultDt string `json:"resultDt"`
}

func wasmValidScript(_ js.Value, args []js.Value) interface{} {

	var plcVars map[string]string
	if err := unmarshalJSON(args[0], &plcVars); err != nil {
		return err.Error()
	}
	resultDt, err := ValidScript(plcVars, args[1].String())
	result := ValidationResult{
		ResultDt: resultDt,
	}
	if err != nil {
		result.Error = err.Error()
	}
	return marshalJSON(result)
}
