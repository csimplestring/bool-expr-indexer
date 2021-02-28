package main

import (
	"container/heap"
	"fmt"
	"testing"
)

func Test_newPostingLists2(t *testing.T) {
	lists := []*PostingList{
		{
			Items: []*PostingItem{
				{CID: 3, Contains: true},
				{CID: 4, Contains: true},
			},
		},
		{
			Items: []*PostingItem{
				{CID: 1, Contains: true},
				{CID: 2, Contains: true},
				{CID: 3, Contains: true},
			},
		},
		{
			Items: []*PostingItem{
				{CID: 3, Contains: false},
				{CID: 4, Contains: true},
			},
		},
	}
	var pCursors []*pCursor
	for _, l := range lists {
		pCursors = append(pCursors, newCursor(l))
	}

	var items []*Item
	pq := make(PriorityQueue, len(lists))
	for i, c := range pCursors {
		item := &Item{value: c, index: i}
		items = append(items, item)
		pq[i] = item
	}

	heap.Init(&pq)

	items[1].value.skipTo(3)
	pq.update(items[1])

	items[0].value.skipTo(4)
	pq.update(items[0])
	items[1].value.skipTo(4)
	pq.update(items[1])
	items[2].value.skipTo(4)
	pq.update(items[2])

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%d %+v", item.value.cur, (*item.value.ref).Items)
	}
}
