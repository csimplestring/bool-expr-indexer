package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type intEle int64

func (i intEle) SortID() int64 { return int64(i) }

func Test_sortedSet(t *testing.T) {
	s := NewSortedSet()

	assert.True(t, s.Add(intEle(3)))
	assert.True(t, s.Add(intEle(1)))
	assert.True(t, s.Add(intEle(2)))

	assert.False(t, s.Add(intEle(3)))
	assert.False(t, s.Add(intEle(1)))
	assert.False(t, s.Add(intEle(2)))

	r := s.ToSlice()
	assert.Equal(t, 3, len(r))
	assert.Equal(t, 3, cap(r))
	for i, e := range r {
		assert.Equal(t, int64(i+1), e.SortID())
	}
}
