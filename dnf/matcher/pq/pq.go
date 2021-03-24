package pq

import "github.com/csimplestring/deheap"

type Item interface {
	Value() interface{}
	Priority() int
	UUID() uint64
}

type IntItem struct {
	Val   int
	Prior int
}

func (i *IntItem) Value() interface{} {
	return i.Val
}

func (i *IntItem) Priority() int {
	return i.Prior
}

func (i *IntItem) UUID() uint64 {
	return uint64(i.Val)
}

type MinMaxPriorityQueue interface {
	Push(Item)
	PeekMin() Item
	PeekMax() Item
	PopMin() Item
	PopMax() Item
	Update(Item)
	Len() int
}

type indexedItem struct {
	value Item
	index int
}

type itemDeheap []*indexedItem

func (h itemDeheap) Len() int { return len(h) }

func (h itemDeheap) Less(i, j int) bool {
	return h[i].value.Priority() < h[j].value.Priority()
}

func (h itemDeheap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *itemDeheap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	n := len(*h)
	item := x.(*indexedItem)
	item.index = n
	*h = append(*h, item)
}

func (h *itemDeheap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	x.index = -1
	*h = old[0 : n-1]
	return x
}

type bpq struct {
	heap  itemDeheap
	idx   map[uint64]int
	bound int
}

func New(bound int) MinMaxPriorityQueue {
	return &bpq{
		heap:  make(itemDeheap, 0, bound),
		idx:   make(map[uint64]int),
		bound: bound,
	}
}

func (b *bpq) Len() int {
	return b.heap.Len()
}

func (b *bpq) add(item Item) {
	it := &indexedItem{
		value: item,
	}
	deheap.Push(&b.heap, it)
	b.idx[item.UUID()] = it.index

	// check bound
	for b.heap.Len() > b.bound {
		ele := deheap.Pop(&b.heap).(*indexedItem)
		delete(b.idx, ele.value.UUID())
	}
}

func (b *bpq) Push(item Item) {

	_, ok := b.idx[item.UUID()]
	if ok {
		b.Update(item)
		return
	}

	b.add(item)
}

func (b *bpq) PeekMin() Item {
	if b.heap.Len() == 0 {
		return nil
	}

	return b.heap[0].value
}

func (b *bpq) PeekMax() Item {
	if b.heap.Len() == 0 {
		return nil
	}

	if b.heap.Len() == 1 {
		return b.heap[0].value
	}
	if b.heap.Len() == 2 {
		return b.heap[1].value
	}

	return max(b.heap[1].value, b.heap[2].value)
}

func (b *bpq) PopMin() Item {
	if b.heap.Len() == 0 {
		return nil
	}

	v := deheap.Pop(&b.heap).(*indexedItem)
	delete(b.idx, v.value.UUID())
	return v.value
}

func (b *bpq) PopMax() Item {
	if b.heap.Len() == 0 {
		return nil
	}

	v := deheap.PopMax(&b.heap).(*indexedItem)
	delete(b.idx, v.value.UUID())
	return v.value
}

func (b *bpq) Update(item Item) {
	if item == nil {
		return
	}

	i, ok := b.idx[item.UUID()]
	if !ok {
		return
	}

	v := deheap.Remove(&b.heap, i).(*indexedItem)
	delete(b.idx, v.value.UUID())

	b.add(item)
}

func max(a Item, b Item) Item {
	if a.Priority() > b.Priority() {
		return a
	}
	return b
}
