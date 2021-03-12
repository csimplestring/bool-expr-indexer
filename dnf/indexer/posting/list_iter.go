package posting

import "sort"

// ListIter ...
type ListIter interface {
	Current() (EntryInt32, bool)
	SkipTo(ID uint32)
}

type listIter struct {
	entries []EntryInt32
	cur     int
}

// NewList creates a read only List
func NewList(entries []EntryInt32) ListIter {
	return &listIter{
		entries: entries,
		cur:     0,
	}
}

func (p *listIter) Current() (EntryInt32, bool) {
	if p.cur >= len(p.entries) {
		return EOL(), false
	}

	return p.entries[p.cur], true
}

func (p *listIter) SkipTo(ID uint32) {
	n := len(p.entries)
	// since p.ref.Items is already sorted in asc order, we do binary search: find the smallest-ID >= ID
	p.cur = sort.Search(n, func(i int) bool { return p.entries[i].CID() >= ID })
}

// Lists ...
type Lists interface {
	SortByCurrent()
	Get(i int) ListIter
	Len() int
}

// NewLists ...
func NewLists(l []ListIter) Lists {
	return lists(l)
}

type lists []ListIter

func (l lists) SortByCurrent() {
	sort.Slice(l[:], func(i, j int) bool {
		a, _ := l[i].Current()
		b, _ := l[j].Current()

		if a.CID() != b.CID() {
			return a.CID() < b.CID()
		}

		return !a.Contains() && b.Contains()
	})
}

func (l lists) Get(i int) ListIter {
	return l[i]
}

func (l lists) Len() int {
	return len(l)
}
