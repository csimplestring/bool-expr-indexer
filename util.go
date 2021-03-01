package main

import "sync"

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

type int64Set struct {
	m map[int64]bool
	sync.RWMutex
}

func newInt64Set() *int64Set {
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
