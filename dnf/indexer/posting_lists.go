package indexer

import (
	"math"
	"sort"
)

// eolItem means end-of-list item, used as the end of a posting list.
var eolItem *PostingEntry = &PostingEntry{
	score:    0,
	CID:      math.MaxInt64,
	Contains: true,
}

// PostingEntry store conjunction-id, belongs-to flag, serving as inverted index pointing to Conjunction
type PostingEntry struct {
	CID      int
	Contains bool
	score    int
}

// PostingList is a list of PostingItem
type PostingList []*PostingEntry

func (p PostingList) sort() {
	sort.Slice(p[:], func(i, j int) bool {

		if p[i].CID != p[j].CID {
			return p[i].CID < p[j].CID
		}

		return !p[i].Contains && p[j].Contains
	})
}

type pCursor struct {
	ref PostingList
	cur int
}

func newCursor(ref PostingList) *pCursor {
	return &pCursor{
		ref: ref,
		cur: 0,
	}
}

func (p *pCursor) current() *PostingEntry {
	if p.cur >= len(p.ref) {
		return eolItem
	}

	return p.ref[p.cur]
}

func (p *pCursor) skipTo(ID int) {
	n := len(p.ref)
	// since p.ref.Items is already sorted in asc order, we do binary search: find the smallest-ID >= ID
	p.cur = sort.Search(n, func(i int) bool { return p.ref[i].CID >= ID })
}

type postingLists []*pCursor

func newPostingLists(l []PostingList) postingLists {
	c := make([]*pCursor, len(l))

	for i, v := range l {
		c[i] = newCursor(v)
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
