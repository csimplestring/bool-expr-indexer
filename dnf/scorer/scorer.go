package scorer

import (
	"fmt"
	"sync"
)

type Scorer interface {
	SetUB(name, value string, ub int)
	GetUB(name, value string) int
}

func NewMapScorer() Scorer {
	return &staticMapScorer{
		m: make(map[string]int),
	}
}

type staticMapScorer struct {
	sync.RWMutex
	m map[string]int
}

func (s *staticMapScorer) key(name, value string) string {
	return fmt.Sprintf("%s:%s", name, value)
}

func (s *staticMapScorer) SetUB(name string, value string, ub int) {
	key := s.key(name, value)

	s.Lock()
	s.m[key] = ub
	s.Unlock()
}

func (s *staticMapScorer) GetUB(name string, value string) int {
	key := s.key(name, value)

	s.RLock()
	defer s.RUnlock()

	return s.m[key]
}
