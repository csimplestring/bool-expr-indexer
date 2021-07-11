package matcher

import (
	"github.com/csimplestring/bool-expr-indexer/api/dnf/indexer"
	"github.com/csimplestring/bool-expr-indexer/api/dnf/indexer/posting"
	"github.com/csimplestring/bool-expr-indexer/api/dnf/scorer"
)

type plistIter struct {
	name  string
	value string
	ub    int
	ref   posting.List
	cur   int
}

func newIterator(ref posting.List) *plistIter {
	return &plistIter{
		ref: ref,
		cur: 0,
	}
}

func newScoredIterator(name, value string, ub int, ref posting.List) *plistIter {
	return &plistIter{
		name:  name,
		value: value,
		ub:    ub,
		ref:   ref,
		cur:   0,
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

func newPostingLists(
	records []*indexer.Record,
	scorer scorer.Scorer,
) postingLists {

	c := make([]*plistIter, len(records))

	for i, v := range records {
		if scorer != nil {
			ub := scorer.GetUB(v.Key, v.Value)
			c[i] = newScoredIterator(v.Key, v.Value, ub, v.PostingList)
		} else {
			c[i] = newIterator(v.PostingList)
		}
	}
	return c
}

func (p postingLists) Len() int {
	return len(p)
}

func (p postingLists) sortByCurrent() {
	// we implement the insertion sort by own because:
	// 1. the size of postingLists is usually small and the changes of position happens not frequently
	// 2. the built-in sort.Sort function takes much time and extra allocation happens, benchmark shows 5x times slower
	p.insertionSort()
}

func (p postingLists) insertionSort() {
	var n = len(p)
	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			a := p[j-1].current()
			b := p[j].current()

			if a.CID() > b.CID() || (a.CID() == b.CID()) && a.Contains() && !b.Contains() {
				p[j-1], p[j] = p[j], p[j-1]
			}
			j = j - 1
		}
	}
}
