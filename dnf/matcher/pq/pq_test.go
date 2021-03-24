package pq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type intItem struct {
	v int
	p int
}

func (i *intItem) Value() interface{} {
	return i.v
}

func (i *intItem) Priority() int {
	return i.p
}

func (i *intItem) UUID() uint64 {
	return uint64(i.v)
}

func TestNew(t *testing.T) {

	pq := New(5)

	num := []int{5, 1, 4, 3, 2}
	items := make([]*intItem, 5)
	for i := 0; i < 5; i++ {
		items[i] = &intItem{
			v: num[i],
			p: num[i],
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
	pq.Push(&intItem{v: 6, p: 6})
	assert.Equal(t, 2, pq.PeekMin().Value())
	assert.Equal(t, 6, pq.PeekMax().Value())

	pq.Push(&intItem{v: 5, p: 100})
	assert.Equal(t, 2, pq.PeekMin().Value())
	assert.Equal(t, 5, pq.PeekMax().Value())

	pq.Push(&intItem{v: 7, p: 7})
	pq.Push(&intItem{v: 8, p: 8})
	assert.Equal(t, 5, pq.PeekMax().Value())
	assert.Equal(t, 4, pq.PeekMin().Value())

	assert.Equal(t, 5, pq.Len())
}
