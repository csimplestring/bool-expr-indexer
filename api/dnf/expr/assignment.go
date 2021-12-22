package expr

import (
	"errors"
	"sync"
)

// Label is a simple k/v pair: like <age:30>
type Label struct {
	Name   string
	Value  string
	Weight int
}

// Assignment is a slice of Label, equals to 'assignment S' in the paper
type Assignment []Label

// mapPool is the pool to recycle the map
var mapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]bool)
	},
}

// ErrorDuplicateLabelName indicates there is duplicated label name in an assignment
var ErrorDuplicateLabelName = errors.New("contains duplicated label name in assignment")

// ValidateAssignment does sanity check on assignment
func ValidateAssignment(a Assignment) error {
	m := mapPool.Get().(map[string]bool)
	for k := range m {
		delete(m, k)
	}
	defer mapPool.Put(m)

	for _, l := range a {
		if _, exist := m[l.Name]; exist {
			return ErrorDuplicateLabelName
		}
		m[l.Name] = true
	}
	return nil
}
