package container

// Iterator ...
type Iterator interface {
	Next() interface{}
}

// emptyIterator
type emptyIterator struct{}

func (i emptyIterator) Next() interface{} { return nil }

var EmptyIterator = emptyIterator{}

// For ...
type IteratorVisit func(interface{}) (over bool)

func For(iter Iterator, fn IteratorVisit) {
	for {
		e := iter.Next()
		if e == nil || fn(e) {
			break
		}
	}
}
