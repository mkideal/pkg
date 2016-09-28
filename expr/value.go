package expr

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrFailToParseInteger    = errors.New("fail to parse integer")
	ErrFailToParseFloat      = errors.New("fail to parse float")
	ErrNotAnInteger          = errors.New("not an integer")
	ErrNotAFloat             = errors.New("not a float")
	ErrUnsupportedType       = errors.New("unsupported type")
	ErrTypeMismatchForOp     = errors.New("type mismatch for operater")
	ErrDivideZero            = errors.New("divide zero")
	ErrPowOfZero             = errors.New("power of zero")
	ErrComparedTypesMismatch = errors.New("compared types mismatch")
)

type Kind int

const (
	KindInvalid Kind = iota
	KindInt
	KindFloat
	KindString
)

var (
	nilValue = Var{kind: KindInvalid}

	varZero  = Var{kind: KindInt, rawValue: "0"}
	varTrue  = Var{kind: KindInt, intValue: 1, rawValue: "true"}
	varFalse = Var{kind: KindInt, intValue: 0, rawValue: "false"}
)

func Nil() Var   { return nilValue }
func Zero() Var  { return varZero }
func True() Var  { return varTrue }
func False() Var { return varFalse }

func Bool(ok bool) Var {
	if ok {
		return True()
	}
	return False()
}

func Int(i int64) Var     { return Var{kind: KindInt, intValue: i, rawValue: strconv.FormatInt(i, 10)} }
func Float(f float64) Var { return Var{kind: KindFloat, floatValue: f, rawValue: fmt.Sprintf("%f", f)} }
func String(s string) Var { return Var{kind: KindString, rawValue: s} }

type Var struct {
	name       string
	kind       Kind
	rawValue   string
	intValue   int64
	floatValue float64
}

func NewVar(name string, kind Kind) Var {
	return Var{name: name, kind: kind}
}

var numberBases = []int{10, 16, 8}

func (v Var) Set(s string) error {
	switch v.kind {
	case KindString:
		// donothing
	case KindInt:
		isInt := false
		for _, base := range numberBases {
			if i, err := strconv.ParseInt(s, base, 64); err == nil {
				v.intValue = i
				isInt = true
				break
			}
		}
		if !isInt {
			return ErrFailToParseInteger
		}
	case KindFloat:
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			v.floatValue = f
		} else {
			return ErrFailToParseFloat
		}
	default:
		return ErrUnsupportedType
	}
	v.rawValue = s
	return nil
}

func (v Var) Name() string   { return v.name }
func (v Var) Kind() Kind     { return v.kind }
func (v Var) String() string { return v.rawValue }
func (v Var) Int() int64 {
	if v.kind == KindFloat {
		return int64(v.floatValue)
	}
	return v.intValue
}
func (v Var) Float() float64 {
	if v.kind == KindInt {
		return float64(v.intValue)
	}
	return v.floatValue
}
func (v Var) Bool() bool {
	switch v.kind {
	case KindString:
		return v.rawValue != ""
	case KindInt:
		return v.intValue != 0
	case KindFloat:
		return v.floatValue != 0
	}
	return false
}

func (v Var) Add(v2 Var) (Var, error) {
	if v.kind == KindString && v2.kind == KindString {
		return stringAdd(v, v2), nil
	}
	return binaryOp(v, v2, intAdd, floatAdd)
}

func (v Var) Sub(v2 Var) (Var, error) { return binaryOp(v, v2, intSub, floatSub) }
func (v Var) Mul(v2 Var) (Var, error) { return binaryOp(v, v2, intMul, floatMul) }
func (v Var) Quo(v2 Var) (Var, error) { return binaryOp(v, v2, intQuo, floatQuo) }
func (v Var) Rem(v2 Var) (Var, error) { return binaryOp(v, v2, intRem, floatRem) }
func (v Var) Pow(v2 Var) (Var, error) { return binaryOp(v, v2, intPow, floatPow) }
func (v Var) And(v2 Var) Var          { return Bool(v.Bool() && v2.Bool()) }
func (v Var) Or(v2 Var) Var           { return Bool(v.Bool() || v2.Bool()) }
func (v Var) Not() Var                { return Bool(!v.Bool()) }
func (v Var) Eq(v2 Var) (Var, error)  { return compare(v, v2, stringEq, intEq, floatEq) }

func (v Var) Neq(v2 Var) (Var, error) {
	result, err := v.Eq(v2)
	if err == nil {
		result = result.Not()
	}
	return result, err
}

func (v Var) Gt(v2 Var) (Var, error) { return compare(v, v2, stringGt, intGt, floatGt) }
func (v Var) Ge(v2 Var) (Var, error) { return compare(v, v2, stringGe, intGe, floatGe) }
func (v Var) Lt(v2 Var) (Var, error) { return v2.Gt(v) }
func (v Var) Le(v2 Var) (Var, error) { return v2.Ge(v) }

func (v Var) Contains(v2 Var) Var {
	if v.kind == KindString && v2.kind == KindString {
		return Bool(strings.Contains(v.rawValue, v2.rawValue))
	}
	return False()
}
