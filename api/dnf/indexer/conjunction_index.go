package indexer

import (
	"strconv"

	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"

	cmap "github.com/orcaman/concurrent-map"
)

type ConjunctionIndex interface {
	Set(ID int, c *expr.Conjunction)
	Get(ID int) (*expr.Conjunction, bool)
}

type conjunctionMap struct {
	m cmap.ConcurrentMap
}

func (c *conjunctionMap) Set(ID int, conjunction *expr.Conjunction) {
	c.m.Set(strconv.Itoa(ID), conjunction)
}

func (c *conjunctionMap) Get(ID int) (*expr.Conjunction, bool) {
	v, exist := c.m.Get(strconv.Itoa(ID))
	if !exist {
		return nil, false
	}

	return v.(*expr.Conjunction), true
}
