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

func (p *plistIter) skipTo(ID uint32) {
	n := len(p.ref)
	// since p.ref.Items is already sorted in asc order, we do search: find the smallest-ID >= ID
	// the binary search is not used
	i := p.cur
	for i < n && p.ref[i].CID() < ID {
		i++
	}
	p.cur = i
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

func (p postingLists) Len() int {
	return len(p)
}

func (p postingLists) sortByCurrent() {
	// we implement the selective sort by own because:
	// 1. the size of postingLists is usually small and the changes of position happens not frequently
	// 2. the built-in sort.Sort function takes much time and extra allocation happens, benchmark shows 5x times slower
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
