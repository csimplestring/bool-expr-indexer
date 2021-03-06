package dnf

import "sort"

// PostingEntry store conjunction-id, belongs-to flag, serving as inverted index pointing to Conjunction
type PostingEntry struct {
	CID      int
	Contains bool
	score    int
}

// PostingList is a list of PostingItem
type PostingList struct {
	Items []*PostingEntry
}

func newPostingList() *PostingList {
	return &PostingList{}
}

func (p *PostingList) append(item *PostingEntry) {
	p.Items = append(p.Items, item)
}

func (p *PostingList) sort() {
	sort.Slice(p.Items[:], func(i, j int) bool {

		if p.Items[i].CID != p.Items[j].CID {
			return p.Items[i].CID < p.Items[j].CID
		}

		return !p.Items[i].Contains && p.Items[j].Contains
	})
}

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

		if a.CID != b.CID {
			return a.CID < b.CID
		}

		return !a.Contains && b.Contains
	})
}

type pCursor struct {
	ref *PostingList
	cur int
}

func newCursor(ref *PostingList) *pCursor {
	return &pCursor{
		ref: ref,
		cur: 0,
	}
}

func (p *pCursor) current() *PostingEntry {
	if p.cur >= len(p.ref.Items) {
		return eolItem
	}

	return p.ref.Items[p.cur]
}

func (p *pCursor) skipTo(ID int) {
	n := len(p.ref.Items)
	// since p.ref.Items is already sorted in asc order, we do binary search: find the smallest-ID >= ID
	p.cur = sort.Search(n, func(i int) bool { return p.ref.Items[i].CID >= ID })
}
