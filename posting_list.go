package main

// import "sort"

// // PostingList is a list of PostingItem
// type PostingList struct {
// 	Items []*PostingItem
// 	cur   int
// }

// func newPostingList() *PostingList {
// 	return &PostingList{
// 		cur: 0,
// 	}
// }

// func (p *PostingList) append(item *PostingItem) {
// 	p.Items = append(p.Items, item)
// }

// func (p *PostingList) sort() {
// 	sort.Slice(p.Items[:], func(i, j int) bool {

// 		if p.Items[i].CID != p.Items[j].CID {
// 			return p.Items[i].CID < p.Items[j].CID
// 		}

// 		return !p.Items[i].Contains && p.Items[j].Contains
// 	})
// }

// func (p *PostingList) rewind() {
// 	p.cur = 0
// }

// func (p *PostingList) current() *PostingItem {
// 	if p.cur >= len(p.Items) {
// 		return eolItem
// 	}

// 	return p.Items[p.cur]
// }

// func (p *PostingList) skipTo(ID int64) {
// 	i := p.cur
// 	for i < len(p.Items) && p.Items[i].CID < ID {
// 		i++
// 	}

// 	p.cur = i
// }

// type postingLists struct {
// 	c []*PostingList
// }

// func newPostingLists(l []*PostingList) *postingLists {
// 	return &postingLists{
// 		c: l,
// 	}
// }

// func (p *postingLists) len() int {
// 	return len(p.c)
// }

// func (p *postingLists) sortByCurrent() {
// 	sort.Slice(p.c[:], func(i, j int) bool {
// 		a := p.c[i].current()
// 		b := p.c[j].current()

// 		// fix, eol item must be in the tail
// 		if a == eolItem && b != eolItem {
// 			return false
// 		}
// 		if a != eolItem && b == eolItem {
// 			return true
// 		}

// 		if a.CID != b.CID {
// 			return a.CID < b.CID
// 		}

// 		return !a.Contains && b.Contains
// 	})
// }
