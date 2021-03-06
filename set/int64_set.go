package set

import "sync"

// IntSet is the set for int64 element.
type IntSet interface {
	Add(x int)
	ToSlice() []int
}

type intSet struct {
	m map[int]bool
	sync.RWMutex
}

// IntHashSet creates a new set for int64 element.
func IntHashSet() IntSet {
	return &intSet{
		m: make(map[int]bool),
	}
}

func (s *intSet) Add(x int) {
	s.Lock()
	s.m[x] = true
	s.Unlock()
}

func (s *intSet) ToSlice() []int {
	s.RLock()
	defer s.RUnlock()

	r := make([]int, len(s.m))
	i := 0
	for k := range s.m {
		r[i] = k
		i++
	}
	return r
}
