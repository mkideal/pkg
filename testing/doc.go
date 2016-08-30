package testing

import (
	"sync"
)

// (NOTE): Stub structure
type common struct {
	mu     sync.Mutex
	output []byte
}

// testing.common extension
