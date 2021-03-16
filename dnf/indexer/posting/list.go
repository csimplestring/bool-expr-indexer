package posting

import "sort"

// List is a list of PostingItem
type List []EntryInt32

func (l List) Sort() {
	sort.Slice(l[:], func(i, j int) bool {

		if l[i].CID() != l[j].CID() {
			return l[i].CID() < l[j].CID()
		}

		return !l[i].Contains() && l[j].Contains()
	})
}
