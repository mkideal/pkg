package expr

import (
	"fmt"
	"go/ast"
	"go/token"
	"math"
	"strconv"
)

//-------------
// + - * / % ^
//-------------

// +
func add(x, y float64) (float64, error) {
	return x + y, nil
}

// -
func sub(x, y float64) (float64, error) {
	return x - y, nil
}

// *
func mul(x, y float64) (float64, error) {
	return x * y, nil
}

// /
func quo(x, y float64) (float64, error) {
	if y == 0 {
		return 0, fmt.Errorf("divided 0")
	}
	return x / y, nil
}

// %
func rem(x, y float64) (float64, error) {
	return math.Mod(x, y), nil
}

// ^
func pow(x, y float64) (float64, error) {
	return math.Pow(x, y), nil
}

// eval the expression
func eval(e *Expr, getter VarGetter, node ast.Expr) (float64, error) {
	switch n := node.(type) {
	case *ast.Ident:
		if getter == nil {
			return e.pool.onVarMissing(n.Name)
		}
		val, ok := getter.GetVar(n.Name)
		if !ok {
			return e.pool.onVarMissing(n.Name)
		}
		return val, nil

	case *ast.BasicLit:
		return strconv.ParseFloat(n.Value, 64)

	case *ast.ParenExpr:
		return eval(e, getter, n.X)

	case *ast.CallExpr:
		args := make([]float64, 0, len(n.Args))
		for _, arg := range n.Args {
			if val, err := eval(e, getter, arg); err != nil {
				return 0, err
			} else {
				args = append(args, val)
			}
		}
		if fnIdent, ok := n.Fun.(*ast.Ident); ok {
			if fn, ok := e.pool.fn(fnIdent.Name); !ok {
				return 0, fmt.Errorf("undefined function `%v`", fnIdent.Name)
			} else {
				return fn(args...)
			}
		}
		return 0, fmt.Errorf("unexpected func type: %T", n.Fun)

	case *ast.UnaryExpr:
		return eval(e, getter, n.X)

	case *ast.BinaryExpr:
		x, err := eval(e, getter, n.X)
		if err != nil {
			return 0, err
		}
		y, err := eval(e, getter, n.Y)
		if err != nil {
			return 0, err
		}
		switch n.Op {
		case token.ADD:
			return add(x, y)
		case token.SUB:
			return sub(x, y)
		case token.MUL:
			return mul(x, y)
		case token.QUO:
			return quo(x, y)
		case token.REM:
			return rem(x, y)
		case token.XOR:
			return pow(x, y)
		default:
			return 0, fmt.Errorf("unexpected binary operator: %v", n.Op)
		}

	default:
		return 0, fmt.Errorf("unexpected node type %T", n)
	}
}
