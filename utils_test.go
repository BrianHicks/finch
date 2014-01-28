package finch

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"time"
)

func TestReverseTaskSlice(t *testing.T) {
	a := NewTask("a", time.Now())
	b := NewTask("b", time.Now())

	orig := []*Task{a, b}
	should := []*Task{b, a}

	assert.Equal(t, should, reverseTaskSlice(orig))
}
