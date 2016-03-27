package expr

import (
	"fmt"
	"math/rand"
	"regexp"
	"sync"
)

type Pool struct {
	locker sync.RWMutex
	pool   map[string]*Expr

	factory map[string]Func
}

func NewPool(factories ...map[string]Func) (*Pool, error) {
	p := new(Pool)
	p.pool = make(map[string]*Expr)
	p.factory = newDefaultFactory()
	for _, factory := range factories {
		if factory == nil {
			continue
		}
		for name, fn := range factory {
			if !validateFuncName(name) {
				return nil, fmt.Errorf("illegal function name `%s`", name)
			}
			p.factory[name] = fn
		}
	}
	return p, nil
}

func (p *Pool) get(s string) (*Expr, bool) {
	p.locker.RLock()
	defer p.locker.RUnlock()
	e, ok := p.pool[s]
	return e, ok && e != nil
}

func (p *Pool) set(s string, e *Expr) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.pool[s] = e
}

func (p *Pool) fn(name string) (Func, bool) {
	fn, ok := p.factory[name]
	return fn, ok
}

// validate function name
var funcNameRegexp = regexp.MustCompile("[a-zA-Z_][a-z-A-Z_0-9]{0,254}")

func validateFuncName(name string) bool {
	return funcNameRegexp.MatchString(name)
}

// default Pool
var defaultPool = func() *Pool {
	p, err := NewPool()
	if err != nil {
		panic(err)
	}
	return p
}()

// default factory
var newDefaultFactory = func() map[string]Func {
	return map[string]Func{
		"min":  minFn,
		"max":  maxFn,
		"rand": randFn,
	}
}

//------------------
// builtin function
//------------------

func minFn(args ...float64) (float64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("missing arguments for function `min`")
	}
	x := args[0]
	for i, size := 1, len(args); i < size; i++ {
		if args[i] < x {
			x = args[i]
		}
	}
	return x, nil
}

func maxFn(args ...float64) (float64, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("missing arguments for function `max`")
	}
	x := args[0]
	for i, size := 1, len(args); i < size; i++ {
		if args[i] > x {
			x = args[i]
		}
	}
	return x, nil
}

func randFn(args ...float64) (float64, error) {
	if len(args) == 0 {
		return float64(rand.Intn(10000)), nil
	}
	if len(args) == 1 {
		x := int(args[0])
		if x <= 0 {
			return 0, fmt.Errorf("bad argument for function `rand`: argument %v <= 0", x)
		}
	}
	if len(args) == 2 {
		x, y := int(args[0]), int(args[1])
		if x > y {
			return 0, fmt.Errorf("bad arguments for function `rand`: first > second")
		}
		return float64(rand.Intn(y-x+1) + x), nil
	}
	return 0, fmt.Errorf("too many arguments for function `rand`: arguments size=%d", len(args))
}
