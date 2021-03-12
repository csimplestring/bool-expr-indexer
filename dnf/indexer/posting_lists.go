package indexer

import (
	"sort"

	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
)

// postingList is a list of PostingItem
type postingList []posting.EntryInt32

func (p postingList) sort() {
	sort.Slice(p[:], func(i, j int) bool {

		if p[i].CID() != p[j].CID() {
			return p[i].CID() < p[j].CID()
		}

		return !p[i].Contains() && p[j].Contains()
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

func (p *plistIter) current() posting.EntryInt32 {
	if p.cur >= len(p.ref) {
		return posting.EOL()
	}

	return p.ref[p.cur]
}

func (p *plistIter) skipTo(ID int) {
	n := len(p.ref)
	// since p.ref.Items is already sorted in asc order, we do binary search: find the smallest-ID >= ID
	p.cur = sort.Search(n, func(i int) bool { return int(p.ref[i].CID()) >= ID })
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

		if a.CID() != b.CID() {
			return a.CID() < b.CID()
		}

		return !a.Contains() && b.Contains()
	})
}
