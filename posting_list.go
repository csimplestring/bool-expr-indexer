package main

import "sort"

type postingLists struct {
	c []*pCursor
}

func newPostingLists(l []*PostingList) *postingLists {
	var c []*pCursor

	for _, v := range l {
		c = append(c, newCursor(v))
	}
	return &postingLists{
		c: c,
	}
}

func (p *postingLists) len() int {
	return len(p.c)
}

func (p *postingLists) sortByCurrent() {
	sort.Slice(p.c[:], func(i, j int) bool {
		a := p.c[i].current()
		b := p.c[j].current()

		// fix, eol item must be in the tail
		if a == eolItem && b != eolItem {
			return false
		}
		if a != eolItem && b == eolItem {
			return true
		}

		if a.CID != b.CID {
			return a.CID < b.CID
		}

		return !a.Contains && b.Contains
	})
}

type pCursor struct {
	ref  *PostingList
	size int
	cur  int
}

func newCursor(ref *PostingList) *pCursor {
	return &pCursor{
		ref:  ref,
		size: len(ref.Items),
		cur:  0,
	}
}

func (p *pCursor) current() *PostingItem {
	if p.cur >= p.size {
		return eolItem
	}

	return p.ref.Items[p.cur]
}

func (p *pCursor) skipTo(ID int64) {
	i := p.cur
	for i < p.size && p.ref.Items[i].CID < ID {
		i++
	}

	p.cur = i
}
