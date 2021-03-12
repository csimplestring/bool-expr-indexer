package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
)

// Indexer shards the Indexer by conjunction size.
type Indexer interface {
	Build() error
	MaxKSize() int
	Add(c *expr.Conjunction)
	Match(assignment expr.Assignment) []int
}
