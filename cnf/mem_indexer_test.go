package cnf

import (
	"testing"

	"github.com/csimplestring/bool-expr-indexer/set"
	"github.com/stretchr/testify/assert"
)

func Test_memoryIndexer_add(t *testing.T) {
	m := &memoryIndexer{
		sortedSetMap: make(map[string][]set.SortedSet),
	}

	m.add("A:1", &PostingEntry{ID: int64(1), DisjunctionID: 0})
	m.add("A:1", &PostingEntry{ID: int64(3), DisjunctionID: 0})
	m.add("A:1", &PostingEntry{ID: int64(2), DisjunctionID: 0})

	m.add("A:1", &PostingEntry{ID: int64(3), DisjunctionID: 1})
	m.add("B:1", &PostingEntry{ID: int64(1), DisjunctionID: 0})

	a1 := m.sortedSetMap["A:1"][0].ToSlice()
	for i, s := range a1 {
		assert.Equal(t, int64(i+1), s.SortID())
	}

	a2 := m.sortedSetMap["A:1"][1].ToSlice()
	assert.Equal(t, int64(3), a2[0].SortID())

	b1 := m.sortedSetMap["B:1"][0].ToSlice()
	assert.Equal(t, int64(1), b1[0].SortID())
}

func Test_memoryIndexer_Build(t *testing.T) {
	s1 := set.NewSortedSet()
	s1.Add(&PostingEntry{ID: int64(1), DisjunctionID: 0})
	s1.Add(&PostingEntry{ID: int64(2), DisjunctionID: 0})
	s1.Add(&PostingEntry{ID: int64(3), DisjunctionID: 0})

	s2 := set.NewSortedSet()
	s2.Add(&PostingEntry{ID: int64(3), DisjunctionID: 1})

	s3 := set.NewSortedSet()
	s3.Add(&PostingEntry{ID: int64(1), DisjunctionID: 0})

	m := &memoryIndexer{
		sortedSetMap: map[string][]set.SortedSet{
			"A:1": {
				s1, s2,
			},
			"B:1": {
				s3,
			},
		},
	}

	m.Build()

	plA := m.pListsMap["A:1"]
	assert.Equal(t, 2, len(plA))
	assert.Equal(t, 3, len(plA[0]))
	assert.Equal(t, 1, len(plA[1]))

	plB := m.pListsMap["B:1"]
	assert.Equal(t, 1, len(plB))
	assert.Equal(t, 1, len(plB[0]))
}
