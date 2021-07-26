package indexer

import (
	"strconv"

	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"

	cmap "github.com/orcaman/concurrent-map"
)

// metadata contains all the meta stats.
type metadata struct {
	forwardIdx *forwardIndex
}

// forwardIndex keeps all the conjunctions.
type forwardIndex struct {
	m cmap.ConcurrentMap
}

// Set sets the ID to conjunction.
func (f *forwardIndex) Set(ID int, conjunction *expr.Conjunction) {
	f.m.Set(strconv.Itoa(ID), conjunction)
}

// Get gets the conjunction by ID.
func (f *forwardIndex) Get(ID int) (*expr.Conjunction, bool) {
	v, exist := f.m.Get(strconv.Itoa(ID))
	if !exist {
		return nil, false
	}

	return v.(*expr.Conjunction), true
}
