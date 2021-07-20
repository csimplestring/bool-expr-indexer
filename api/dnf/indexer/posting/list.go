package posting

import "sort"

// List is a list of Posting Entry
type List []EntryInt32

// Sort sorts l by CID and Contains flag in asc
func (l List) Sort() {
	sort.Slice(l[:], func(i, j int) bool {

		if l[i].CID() != l[j].CID() {
			return l[i].CID() < l[j].CID()
		}

		return !l[i].Contains() && l[j].Contains()
	})
}

// Remove removes ID from l, this function is expensive!
func (l List) Remove(ID int) bool {
	idx := l.IndexOf(ID)
	if idx == -1 {
		return false
	}

	l = append(l[:idx], l[idx+1:]...)
	return true
}

func (l List) IndexOf(ID int) int {
	// binary search in l since l is ascending sorted
	i := sort.Search(len(l), func(i int) bool {
		return l[i].CID() >= uint32(ID)
	})

	if i < len(l) && l[i].CID() == uint32(ID) {
		// x is present
		return i
	}

	return -1
}
