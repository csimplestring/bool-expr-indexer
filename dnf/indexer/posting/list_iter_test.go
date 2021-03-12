package posting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewList(t *testing.T) {

	e1, _ := NewEntryInt32(1, true, 0)
	e2, _ := NewEntryInt32(2, true, 0)
	e3, _ := NewEntryInt32(3, true, 0)

	entries := []EntryInt32{
		e1, e2, e3,
	}

	l0 := NewList([]EntryInt32{})
	c, next := l0.Current()
	assert.Equal(t, EOL, c)
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
	assert.Equal(t, EOL, c)
	assert.False(t, next)
}

func TestNewLists(t *testing.T) {

	e1, _ := NewEntryInt32(1, true, 0)
	e2, _ := NewEntryInt32(2, true, 0)
	e3, _ := NewEntryInt32(3, true, 0)
	e4, _ := NewEntryInt32(3, false, 0)

	l1 := NewList([]EntryInt32{
		e1, e2, e3,
	})

	l2 := NewList([]EntryInt32{
		e3,
	})

	l3 := NewList([]EntryInt32{
		e2,
		e4,
	})

	lists := NewLists([]ListIter{l1, l2, l3})
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
