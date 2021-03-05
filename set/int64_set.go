package set

import "sync"

// IntSet is the set for int64 element.
type IntSet interface {
	Add(x int)
	ToSlice() []int
}

type int64Set struct {
	m map[int]bool
	sync.RWMutex
}

// Int64HashSet creates a new set for int64 element.
func Int64HashSet() IntSet {
	return &int64Set{
		m: make(map[int]bool),
	}
}

func (s *int64Set) Add(x int) {
	s.Lock()
	s.m[x] = true
	s.Unlock()
}

func (s *int64Set) ToSlice() []int {
	r := make([]int, len(s.m))
	i := 0
	for k := range s.m {
		r[i] = k
		i++
	}
	return r
}
