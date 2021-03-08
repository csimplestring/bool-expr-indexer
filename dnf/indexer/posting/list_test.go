package posting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewList(t *testing.T) {

	entries := []Entry{
		{CID: 1, Contains: true},
		{CID: 2, Contains: true},
		{CID: 3, Contains: true},
	}

	l0 := NewList([]Entry{})
	c, next := l0.Current()
	assert.Equal(t, *EOL, c)
	assert.False(t, next)

	l1 := NewList(entries)
	c, next = l1.Current()
	assert.Equal(t, entries[0], c)
	assert.True(t, next)

	l1.SkipTo(3)
	c, next = l1.Current()
	assert.Equal(t, entries[2], c)
	assert.True(t, next)

	l1.SkipTo(4)
	c, next = l1.Current()
	assert.Equal(t, *EOL, c)
	assert.False(t, next)
}

func TestNewLists(t *testing.T) {

	l1 := NewList([]Entry{
		{CID: 1, Contains: true},
		{CID: 2, Contains: true},
		{CID: 3, Contains: true},
	})

	l2 := NewList([]Entry{
		{CID: 3, Contains: true},
	})

	l3 := NewList([]Entry{
		{CID: 2, Contains: true},
		{CID: 3, Contains: false},
	})

	lists := NewLists([]List{l1, l2, l3})
	lists.SortByCurrent()

	assert.Equal(t, l1, lists.Get(0))
	assert.Equal(t, l3, lists.Get(1))
	assert.Equal(t, l2, lists.Get(2))

	l1.SkipTo(3)
	l2.SkipTo(3)
	l3.SkipTo(3)
	lists.SortByCurrent()

	assert.Equal(t, l3, lists.Get(0))
	assert.Equal(t, l1, lists.Get(1))
	assert.Equal(t, l2, lists.Get(2))
}
