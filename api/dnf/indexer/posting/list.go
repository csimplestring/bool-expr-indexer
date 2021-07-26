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
func (l *List) Remove(ID int) bool {
	first := l.indexOf(uint32(ID))
	if first == -1 {
		return false
	}
	last := l.lastIndexOf(uint32(ID))

	*l = append((*l)[:first], (*l)[last+1:]...)
	return true
}

func (nums List) indexOf(target uint32) int {
	l, r := 0, len(nums)

	for l < r {
		m := l + (r-l)/2
		if nums[m].CID() >= target {
			r = m
		} else {
			l = m + 1
		}
	}

	if l < len(nums) && nums[l].CID() == target {
		return l
	} else {
		return -1
	}
}

func (nums List) lastIndexOf(target uint32) int {
	l, r := 0, len(nums)-1

	for l < r {
		m := (l + r + 1) / 2

		if nums[m].CID() <= target {
			l = m
		} else {
			r = m - 1
		}
	}

	if l < len(nums) && nums[l].CID() == target {
		return l
	} else {
		return -1
	}
}
