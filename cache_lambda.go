package edgeexpr

import (
	"fmt"
	"reflect"

	"github.com/expr-lang/expr"
	"github.com/samber/lo"
)

func (c *Cache) Contains(v any, window string) (bool, error) {
	if c == nil {
		return false, nil
	}
	points := c.getPointsInWindow(window)
	if len(points) == 0 {
		return false, nil
	}

	// 使用类型断言检查是否为 float64
	return lo.ContainsBy(points, func(p Point) bool {
		return reflect.DeepEqual(p.Value, v)
	}), nil
}

// lambda: (x) => bool
func (c *Cache) EveryBy(lambda, window string) (bool, error) {
	if c == nil {
		return false, nil
	}
	points := c.getPointsInWindow(window)
	if len(points) == 0 {
		return false, nil
	}

	// make sure points in windows is complete loaded
	if len(points) >= len(c.Points) {
		return false, nil
	}
	env := map[string]any{"x": points[0].Value}
	program, err := expr.Compile(lambda, expr.Env(env))
	if err != nil {
		return false, fmt.Errorf("lambda: %s program with err: %v", lambda, err)
	}
	out, err := expr.Run(program, env)
	if err != nil {
		return false, fmt.Errorf("lambda: %s env with err: %v", lambda, err)
	}
	if _, ok := out.(bool); !ok {
		return false, fmt.Errorf("lambda: %s out is not a bool", lambda)
	}
	return lo.EveryBy(points, func(p Point) bool {
		out, _ := expr.Eval(lambda, map[string]any{"x": p.Value})
		return out.(bool)
	}), nil
}

// lambda: (x) => bool
func (c *Cache) SomeBy(lambda, window string) (bool, error) {
	if c == nil {
		return false, nil
	}
	points := c.getPointsInWindow(window)
	if len(points) == 0 {
		return false, nil
	}

	env := map[string]any{"x": points[0].Value}
	program, err := expr.Compile(lambda, expr.Env(env))
	if err != nil {
		return false, fmt.Errorf("lambda: %s program with err: %v", lambda, err)
	}
	out, err := expr.Run(program, env)
	if err != nil {
		return false, fmt.Errorf("lambda: %s env with err: %v", lambda, err)
	}
	if _, ok := out.(bool); !ok {
		return false, fmt.Errorf("lambda: %s out is not a bool", lambda)
	}
	return lo.SomeBy(points, func(p Point) bool {
		out, _ := expr.Eval(lambda, map[string]any{"x": p.Value})
		return out.(bool)
	}), nil
}

// lambda: (x) => bool
func (c *Cache) NoneBy(lambda, window string) (bool, error) {
	if c == nil {
		return false, nil
	}
	points := c.getPointsInWindow(window)
	if len(points) == 0 {
		return false, nil
	}

	// make sure points in windows is complete loaded
	if len(points) >= len(c.Points) {
		return false, nil
	}
	env := map[string]any{"x": points[0].Value}
	program, err := expr.Compile(lambda, expr.Env(env))
	if err != nil {
		return false, fmt.Errorf("lambda: %s program with err: %v", lambda, err)
	}
	out, err := expr.Run(program, env)
	if err != nil {
		return false, fmt.Errorf("lambda: %s env with err: %v", lambda, err)
	}
	if _, ok := out.(bool); !ok {
		return false, fmt.Errorf("lambda: %s out is not a bool", lambda)
	}
	return lo.NoneBy(points, func(p Point) bool {
		out, _ := expr.Eval(lambda, map[string]any{"x": p.Value})
		return out.(bool)
	}), nil
}
