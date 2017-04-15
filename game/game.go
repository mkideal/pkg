package game

import (
	"errors"
)

var (
	ErrUnexpectedGameover = errors.New("unexpected gameover")
	ErrGameover           = errors.New("gameover")
)

type Player interface {
	Id() int64
	Reset()
	Notify(msg interface{}) error
}
