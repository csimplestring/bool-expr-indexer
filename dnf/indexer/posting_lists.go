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
		return posting.EOL
	}

	return p.ref[p.cur]
}

func (p *plistIter) skipTo(ID int) {
	n := len(p.ref)
	// since p.ref.Items is already sorted in asc order, we do binary search: find the smallest-ID >= ID
	p.cur = search(p.cur, n, func(i int) bool { return int(p.ref[i].CID()) >= ID })
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

func (p postingLists) Less(i, j int) bool {
	a := p[i].current()
	b := p[j].current()

	if a.CID() != b.CID() {
		return a.CID() < b.CID()
	}

	return !a.Contains() && b.Contains()
}

func (p postingLists) Swap(i, j int) {
	p[j], p[i] = p[i], p[j]
}

func (p postingLists) Len() int {
	return len(p)
}

func (p postingLists) sortByCurrent() {
	for i := 0; i < len(p)-1; i++ {
		min, index := p[i], i
		for j := i + 1; j < len(p); j++ {
			if p[j].current().CID() < min.current().CID() || (p[j].current().CID() == min.current().CID() && !p[j].current().Contains() && min.current().Contains()) {
				min, index = p[j], j
			}

		}
		if i != index {
			p[i], p[index] = p[index], p[i]
		}
	}
}
