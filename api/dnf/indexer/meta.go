package indexer

import (
	"strconv"

	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"

	cmap "github.com/orcaman/concurrent-map"
)

type metadata struct {
	forwardIdx *forwardIndex
}

type forwardIndex struct {
	m cmap.ConcurrentMap
}

func (f *forwardIndex) Set(ID int, conjunction *expr.Conjunction) {
	f.m.Set(strconv.Itoa(ID), conjunction)
}

func (f *forwardIndex) Get(ID int) (*expr.Conjunction, bool) {
	v, exist := f.m.Get(strconv.Itoa(ID))
	if !exist {
		return nil, false
	}

	return v.(*expr.Conjunction), true
}
