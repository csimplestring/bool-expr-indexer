package main

import "sort"

// Matcher finds the matched conjunction ids
type Matcher interface {
	Match(*kIndexTable, Labels) []int64
}

type matcher struct {
}

func (m *matcher) Match(k *kIndexTable, labels Labels) []int64 {
	var results []int64

	n := min(len(labels), k.maxKSize)

	for i := n; i >= 0; i-- {
		pLists := newPostingLists(k.GetPostingLists(i, labels))

		K := i
		if K == 0 {
			K = 1
		}
		if pLists.len() < K {
			continue
		}

		pLists.sortByCurrent()
		for pLists.c[K-1].current() != eolItem {
			var nextID int64

			if pLists.c[0].current().CID == pLists.c[K-1].current().CID {

				if pLists.c[0].current().Contains == false {
					rejectID := pLists.c[0].current().CID
					for L := K; L <= pLists.len()-1; L++ {
						if pLists.c[L].current().CID == rejectID {
							pLists.c[L].skipTo(rejectID + 1)
						} else {
							break
						}
					}

				} else {
					results = append(results, pLists.c[K-1].current().CID)
				}

				nextID = pLists.c[K-1].current().CID + 1
			} else {
				nextID = pLists.c[K-1].current().CID

			}

			for L := 0; L <= K-1; L++ {
				pLists.c[L].skipTo(nextID)
			}
			pLists.sortByCurrent()
		}

	}

	return results
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
