package container

type Container interface {
	Len() int
	Iter() Iterator
	Contains(interface{}) bool
}

type Array interface {
	Container
}

type Map interface {
	Container
}

// Iterator ...
type Iterator interface {
	Next() (k, v interface{})
}

// emptyIterator
type emptyIterator struct{}

func (i emptyIterator) Next() interface{} { return nil }

var EmptyIterator = emptyIterator{}
