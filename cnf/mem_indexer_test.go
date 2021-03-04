package cnf

import (
	"testing"

	"github.com/csimplestring/bool-expr-indexer/set"
	"github.com/stretchr/testify/assert"
)

func Test_memoryIndexer_add(t *testing.T) {
	m := &memoryIndexer{
		s: make(map[string][]set.SortedSet),
	}

	m.add("A:1", &PostingEntry{ID: int64(1), DisjunctionID: 0})
	m.add("A:1", &PostingEntry{ID: int64(3), DisjunctionID: 0})
	m.add("A:1", &PostingEntry{ID: int64(2), DisjunctionID: 0})

	m.add("A:1", &PostingEntry{ID: int64(3), DisjunctionID: 1})
	m.add("B:1", &PostingEntry{ID: int64(1), DisjunctionID: 0})

	a1 := m.s["A:1"][0].ToSlice()
	for i, s := range a1 {
		assert.Equal(t, int64(i+1), s.SortID())
	}

	a2 := m.s["A:1"][1].ToSlice()
	assert.Equal(t, int64(3), a2[0].SortID())

	b1 := m.s["B:1"][0].ToSlice()
	assert.Equal(t, int64(1), b1[0].SortID())
}
