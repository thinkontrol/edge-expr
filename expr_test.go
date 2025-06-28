package edgeexpr

import (
	"testing"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

func TestCacheExpr(t *testing.T) {
	env := map[string]any{
		"temperature": NewCache[float64](time.Minute),
		"status":      NewCache[bool](time.Minute),
		"message":     NewCache[string](time.Minute),
		"data":        NewCache[[]byte](time.Minute),
	}

	expressions := []string{
		`temperature.Value() * 0.5`,
		`status.Value()`,
		`message.Value() + " is ready"`,
		`data.Value()`,
		`temperature.Len()`,
		`status.Len()`,
		`message.Len()`,
		`data.Len()`,
		`temperature.MA('20s')`,
		`status.MA('20s')`,
		`message.MA('20s')`,
		`data.MA('20s')`,
		`temperature.StdDev('20s')`,
		`status.StdDev('20s')`,
		`message.StdDev('20s')`,
		`data.StdDev('20s')`,
		`temperature.PctChange()`,
		`status.PctChange()`,
		`message.PctChange()`,
		`data.PctChange()`,
		`temperature.Diff()`,
		`status.Diff()`,
		`message.Diff()`,
		`data.Diff()`,
		`temperature.PctChangeExceeds(10)`,
		`status.PctChangeExceeds(10)`,
		`message.PctChangeExceeds(10)`,
		`data.PctChangeExceeds(10)`,
		`temperature.DiffExceeds(10)`,
		`status.DiffExceeds(10)`,
		`message.DiffExceeds(10)`,
		`data.DiffExceeds(10)`,
		`temperature.Changed()`,
		`status.Changed()`,
		`message.Changed()`,
		`data.Changed()`,
		`temperature.PctChangeSince('10s')`,
		`status.PctChangeSince('10s')`,
		`message.PctChangeSince('10s')`,
		`data.PctChangeSince('10s')`,
		`temperature.DiffSince('10s')`,
		`status.DiffSince('10s')`,
		`message.DiffSince('10s')`,
		`data.DiffSince('10s')`,
		`temperature.Count('10s')`,
		`status.Count('10s')`,
		`message.Count('10s')`,
		`data.Count('10s')`,
		`temperature.Rising()`,
		`status.Rising()`,
		`message.Rising()`,
		`data.Rising()`,
		`temperature.Falling()`,
		`status.Falling()`,
		`message.Falling()`,
		`data.Falling()`,
		`temperature.RC('60s')`,
		`status.RC('60s')`,
		`message.RC('60s')`,
		`data.RC('60s')`,
		`temperature.FC('60s')`,
		`status.FC('60s')`,
		`message.FC('60s')`,
		`data.FC('60s')`,
		`temperature.Bit(12)`,
		`status.Bit(12)`,
		`message.Bit(12)`,
		`data.Bit(12)`,
		`temperature.ByteBit(0,8)`,
		`status.ByteBit(0,8)`,
		`message.ByteBit(0,8)`,
		`data.ByteBit(0,8)`,
		`temperature.ByteBit(0,4)`,
		`status.ByteBit(0,4)`,
		`message.ByteBit(0,4)`,
		`data.ByteBit(0,4)`,
		`temperature.ByteBit(2,4)`,
		`status.ByteBit(2,4)`,
		`message.ByteBit(2,4)`,
		`data.ByteBit(2,4)`,
	}

	var programs []*vm.Program

	for _, exprStr := range expressions {
		program, err := expr.Compile(exprStr, expr.Env(env))
		if err != nil {
			t.Errorf("failed to compile expression %q: %v", exprStr, err)
		} else {
			t.Logf("compiled expression %q successfully", exprStr)
		}
		programs = append(programs, program)
	}

}
