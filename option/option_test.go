package option

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOption(t *testing.T) {
	assert.Equal(t, true, Bool(true))
	assert.Equal(t, false, Bool(false))
	assert.Equal(t, true, Bool(true, true))
	assert.Equal(t, false, Bool(true, false))
	assert.Equal(t, true, Bool(false, true))
	assert.Equal(t, false, Bool(false, false))
}
