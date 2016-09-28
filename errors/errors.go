package errors

import (
	"sync/atomic"

	"github.com/mkideal/pkg/debug"
)

// Error aliases string
type Error string

func (e Error) Error() string { return string(e) }

var enableTrace int32 = 1

// SwitchTrace switchs trace on/off
func SwitchTrace(on bool) {
	if on {
		atomic.StoreInt32(&enableTrace, 1)
	} else {
		atomic.StoreInt32(&enableTrace, 0)
	}
}

// TraceError wraps string and contains created stack infos
type TraceError struct {
	stack string
	text  string
}

func (e TraceError) String() string {
	return e.stack + "\n" + e.text
}

func (e TraceError) Error() string {
	if atomic.LoadInt32(&enableTrace) == 0 {
		return e.text
	}
	return e.String()
}

// Throw throws an error which contains stack infos
func Throw(text string) error {
	stack := string(debug.Stack(2))
	return TraceError{stack: stack, text: text}
}

// WrappedError define an interface which wrap another error
type WrappedError interface {
	error
	Core() error
}

var _ = WrappedError(wrappedError{})

// wrappedError implements WrappedError
type wrappedError struct {
	stack string
	err   error
}

func (e wrappedError) String() string {
	return e.stack + "\n" + e.err.Error()
}

func (e wrappedError) Error() string {
	if atomic.LoadInt32(&enableTrace) == 0 {
		return e.err.Error()
	}
	return e.String()
}

func (e wrappedError) Core() error { return e.err }

// Wrap wraps another error
func Wrap(err error) error {
	return wrappedError{stack: string(debug.Stack(2)), err: err}
}

// Core returns wrapped error
func Core(err error) error {
	if err == nil {
		return nil
	}
	const maxWrapLayer = 32
	c := 0
	for c < maxWrapLayer {
		c++
		wrapped, ok := err.(WrappedError)
		if !ok {
			break
		}
		err = wrapped.Core()
	}
	return err
}
