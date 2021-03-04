package set

import "sort"

// Item ...
type Item interface {
	SortID() int64
}

// SortedSet ...
type SortedSet interface {
	Add(e Item) bool
	ToSlice() []Item
}

// NewSortedSet ...
func NewSortedSet() SortedSet {
	return &sortedSet{}
}

type sortedSet struct {
	data []Item
}

func (s *sortedSet) ToSlice() []Item {
	r := make([]Item, len(s.data))
	copy(r, s.data)
	return r
}

// Add inserts e to s if it does not exist and returns true
// if e already exists in s, then rejected and returns false
// the s should be maintained in asc order
func (s *sortedSet) Add(e Item) bool {
	if len(s.data) == 0 {
		s.data = append(s.data, e)
		return true
	}

	// find the first element that larger than e.ID
	idx := sort.Search(len(s.data), func(i int) bool { return s.data[i].SortID() >= e.SortID() })

	// e is present
	if idx < len(s.data) && s.data[idx].SortID() == e.SortID() {
		return false
	}

	// e is not present in data, but i is the index where it would be inserted.
	s.data = append(s.data, nil)
	copy(s.data[idx+1:], s.data[idx:])
	s.data[idx] = e

	return true
}
