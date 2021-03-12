package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_search(t *testing.T) {

	a := []int{}
	id := 0
	idx := search(0, 0, func(i int) bool {
		return a[i] >= id
	})
	assert.Equal(t, 0, idx)

	a = []int{1}
	idx = search(0, 1, func(i int) bool {
		return a[i] >= 0
	})
	assert.Equal(t, 0, idx)

	idx = search(0, 1, func(i int) bool {
		return a[i] >= 1
	})
	assert.Equal(t, 0, idx)

	idx = search(0, 1, func(i int) bool {
		return a[i] >= 2
	})
	assert.Equal(t, 1, idx)

	idx = search(1, 1, func(i int) bool {
		return a[i] >= 2
	})
	assert.Equal(t, 1, idx)

	a = []int{1, 3, 5, 6, 7, 8, 11}
	idx = search(0, len(a), func(i int) bool {
		return a[i] >= 2
	})
	assert.Equal(t, 1, idx)

	idx = search(2, len(a), func(i int) bool {
		return a[i] >= 2
	})
	assert.Equal(t, 2, idx)

	idx = search(2, len(a), func(i int) bool {
		return a[i] >= 5
	})
	assert.Equal(t, 2, idx)

	idx = search(2, len(a), func(i int) bool {
		return a[i] >= 8
	})
	assert.Equal(t, 5, idx)

	idx = search(2, len(a), func(i int) bool {
		return a[i] > 8
	})
	assert.Equal(t, 6, idx)

	idx = search(2, len(a), func(i int) bool {
		return a[i] >= 100
	})
	assert.Equal(t, 7, idx)

	idx = search(3, len(a), func(i int) bool {
		return a[i] >= 100
	})
	assert.Equal(t, 7, idx)

	idx = search(7, len(a), func(i int) bool {
		return a[i] >= 100
	})
	assert.Equal(t, 7, idx)
}
