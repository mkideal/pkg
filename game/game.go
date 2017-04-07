package game

type Player interface {
	Id() int64
	Reset()
}

type Command struct {
	Pos  int
	Data interface{}
}
