package pq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	pq := New(5)

	num := []int{5, 1, 4, 3, 2}
	items := make([]*IntItem, 5)
	for i := 0; i < 5; i++ {
		items[i] = &IntItem{
			Val:   num[i],
			Prior: num[i],
		}
	}

	for _, item := range items {
		pq.Push(item)
	}

	assert.Equal(t, 1, pq.PeekMin().Value())
	assert.Equal(t, 5, pq.PeekMax().Value())

	// pop 1
	pq.PopMin()
	assert.Equal(t, 2, pq.PeekMin().Value())

	// pop 5
	pq.PopMax()
	assert.Equal(t, 4, pq.PeekMax().Value())

	// len = 4
	pq.Push(&IntItem{Val: 6, Prior: 6})
	assert.Equal(t, 2, pq.PeekMin().Value())
	assert.Equal(t, 6, pq.PeekMax().Value())

	pq.Push(&IntItem{Val: 5, Prior: 100})
	assert.Equal(t, 2, pq.PeekMin().Value())
	assert.Equal(t, 5, pq.PeekMax().Value())

	pq.Push(&IntItem{Val: 7, Prior: 7})
	pq.Push(&IntItem{Val: 8, Prior: 8})
	assert.Equal(t, 5, pq.PeekMax().Value())
	assert.Equal(t, 4, pq.PeekMin().Value())

	assert.Equal(t, 5, pq.Len())
}
