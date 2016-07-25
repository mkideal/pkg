package expr

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstExpr(t *testing.T) {
	for _, x := range []struct {
		s   string
		val float64
	}{
		{"1", 1},
		{"1+2", 3},
		{"1*2+3", 5},
		{"1*(2+3)", 5},
	} {
		e, err := New(x.s, nil)
		if err != nil {
			t.Errorf("parse %q error: %v", x.s, err)
			continue
		}
		val, err := e.Eval(nil)
		if err != nil {
			t.Errorf("parse %q error: %v", x.s, err)
			continue
		}
		if math.Abs(val-x.val) > 1E-6 {
			t.Errorf("%q want %f, got %f", x.s, x.val, val)
		}
	}
}

func TestVarExpr(t *testing.T) {
	getter := Getter{
		"x": 1.5,
		"y": 2.5,
		"m": 5,
		"n": 2,
	}

	for _, x := range []struct {
		s     string
		val   float64
		isErr bool
	}{
		{"1", 1, false},
		{"x", 1.5, false},
		{"y", 2.5, false},
		{"x+y", 4, false},
		{"x-y", -1, false},
		{"x*y", 3.75, false},
		{"x/y", 0.6, false},
		{"m / n", 2.5, false},
		{"m / n // line comment", 2.5, false},
		{"m / /*multiline\n	comment*/ n", 2.5, false},
		{"m%n", 1, false},
		{"(1+x)*y/(m-n)^2", 4.34027776, false},
		{"(1+x)*y/((m-n)^2)", 0.69444444, false},

		{"min(x,y,m)", 1.5, false},
		{"max(x,y,m)", 5, false},
		{"max(1, x+y)", 4, false},
		{"min(max(1, x+y), 2, x)", 1.5, false},

		{"undefined_func(x,y,m)", 0, true},
		{"max()", 0, true},
		{"m!-!n", 0, true},
		{"m / undefined_var", 0, true},
	} {
		e, err := New(x.s, nil)
		if err != nil {
			if !x.isErr {
				t.Errorf("parse %q error: %v", x.s, err)
			}
			continue
		}
		val, err := e.Eval(getter)
		if err != nil {
			if !x.isErr {
				t.Errorf("parse %q error: %v", x.s, err)
			}
			continue
		}
		if math.Abs(val-x.val) > 1E-6 {
			t.Errorf("%q want %f, got %f", x.s, x.val, val)
		}
	}
}

func TestCustomFactory(t *testing.T) {
	pool, _ := NewPool(map[string]Func{
		"constant": func(...float64) (float64, error) { return 123, nil },
		"sum": func(args ...float64) (float64, error) {
			sum := float64(0)
			for _, arg := range args {
				sum += arg
			}
			return sum, nil
		},
		"average": func(args ...float64) (float64, error) {
			n := len(args)
			if n == 0 {
				return 0, fmt.Errorf("missing arguments for function `%s`", "average")
			}
			sum := float64(0)
			for _, arg := range args {
				sum += arg
			}
			return sum / float64(n), nil
		},
	})
	getter := Getter{
		"x": 1.5,
		"y": 2.5,
	}

	for _, x := range []struct {
		s     string
		val   float64
		isErr bool
	}{
		{"constant()", 123, false},
		{"constant(x)", 123, false},
		{"sum(1,2,3)", 6, false},
		{"sum()", 0, false},
		{"sum(x, y, x)", 5.5, false},
		{"average(x)", 1.5, false},
		{"average(x,y)", 2, false},
		{"average()", 0, true},
	} {
		e, err := New(x.s, pool)
		if err != nil {
			if !x.isErr {
				t.Errorf("parse %q error: %v", x.s, err)
			}
			continue
		}
		val, err := e.Eval(getter)
		if err != nil {
			if !x.isErr {
				t.Errorf("parse %q error: %v", x.s, err)
			}
			continue
		}
		if math.Abs(val-x.val) > 1E-6 {
			t.Errorf("%q want %f, got %f", x.s, x.val, val)
		}
	}
}

func TestOnVarMissing(t *testing.T) {
	defaults := map[string]float64{
		"a": 0,
		"b": 1,
	}
	pool, _ := NewPool()
	pool.SetOnVarMissing(func(varName string) (float64, error) {
		if dft, ok := defaults[varName]; ok {
			return dft, nil
		}
		return DefaultOnVarMissing(varName)
	})
	v, err := Eval("2 / b + a + x", map[string]float64{"x": 1}, pool)
	assert.Nil(t, err)
	assert.Equal(t, float64(3), v)

	v, err = Eval("2 / b + a + undefined", nil, pool)
	assert.NotNil(t, err)
}
