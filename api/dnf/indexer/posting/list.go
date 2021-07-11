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
