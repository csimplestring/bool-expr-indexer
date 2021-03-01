package set

import "sync"

// Int64Set is the set for int64 element.
type Int64Set interface {
	Add(x int64)
	ToSlice() []int64
}

type int64Set struct {
	m map[int64]bool
	sync.RWMutex
}

// Int64HashSet creates a new set for int64 element.
func Int64HashSet() Int64Set {
	return &int64Set{
		m: make(map[int64]bool),
	}
}

func (s *int64Set) Add(x int64) {
	s.Lock()
	s.m[x] = true
	s.Unlock()
}

func (s *int64Set) ToSlice() []int64 {
	r := make([]int64, len(s.m))
	i := 0
	for k := range s.m {
		r[i] = k
		i++
	}
	return r
}
