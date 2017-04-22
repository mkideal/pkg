package game

import (
	"errors"
)

var (
	ErrUnexpectedGameover = errors.New("unexpected gameover")
	ErrGameover           = errors.New("gameover")
	ErrState              = errors.New("error status")
	ErrTurn               = errors.New("error turn")
	ErrCommandRepeated    = errors.New("command repeated")
)

type Player interface {
	ID() int64
	Reset()
	Notify(msg interface{}) error
}
