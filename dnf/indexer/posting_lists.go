package indexer

import (
	"math"
	"sort"
)

// eolItem means end-of-list item, used as the end of a posting list.
var eolItem *postingEntry = &postingEntry{
	score:    0,
	CID:      math.MaxInt64,
	Contains: true,
}

// postingEntry store conjunction-id, belongs-to flag, serving as inverted index pointing to Conjunction
type postingEntry struct {
	CID      int
	Contains bool
	score    int
}

// postingList is a list of PostingItem
type postingList []postingEntry

func (p postingList) sort() {
	sort.Slice(p[:], func(i, j int) bool {

		if p[i].CID != p[j].CID {
			return p[i].CID < p[j].CID
		}

		return !p[i].Contains && p[j].Contains
	})
}

type plistIter struct {
	ref postingList
	cur int
}

func newIterator(ref postingList) *plistIter {
	return &plistIter{
		ref: ref,
		cur: 0,
	}
}

func (p *plistIter) current() postingEntry {
	if p.cur >= len(p.ref) {
		return *eolItem
	}

	return p.ref[p.cur]
}

func (p *plistIter) skipTo(ID int) {
	n := len(p.ref)
	// since p.ref.Items is already sorted in asc order, we do binary search: find the smallest-ID >= ID
	p.cur = sort.Search(n, func(i int) bool { return p.ref[i].CID >= ID })
}

// postingLists is a slice of list iterator
type postingLists []*plistIter

func newPostingLists(l []postingList) postingLists {
	c := make([]*plistIter, len(l))

	for i, v := range l {
		c[i] = newIterator(v)
	}
	return c
}

func (p postingLists) len() int {
	return len(p)
}

func (p postingLists) sortByCurrent() {
	sort.Slice(p[:], func(i, j int) bool {
		a := p[i].current()
		b := p[j].current()

		if a.CID != b.CID {
			return a.CID < b.CID
		}

		return !a.Contains && b.Contains
	})
}
