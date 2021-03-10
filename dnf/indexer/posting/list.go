package posting

import "sort"

// List ...
type List interface {
	Current() (EntryInt32, bool)
	SkipTo(ID uint32)
}

type list struct {
	entries []EntryInt32
	cur     int
}

// NewList creates a read only List
func NewList(entries []EntryInt32) List {
	return &list{
		entries: entries,
		cur:     0,
	}
}

func (p *list) Current() (EntryInt32, bool) {
	if p.cur >= len(p.entries) {
		return EOL, false
	}

	return p.entries[p.cur], true
}

func (p *list) SkipTo(ID uint32) {
	n := len(p.entries)
	// since p.ref.Items is already sorted in asc order, we do binary search: find the smallest-ID >= ID
	p.cur = sort.Search(n, func(i int) bool { return p.entries[i].CID() >= ID })
}

// Lists ...
type Lists interface {
	SortByCurrent()
	Get(i int) List
	Len() int
}

// NewLists ...
func NewLists(l []List) Lists {
	return lists(l)
}

type lists []List

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

func (l lists) Get(i int) List {
	return l[i]
}

func (l lists) Len() int {
	return len(l)
}
