package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
)

// Indexer defines the top level indexer interface
type Indexer interface {
	Build() error
	MaxKSize() int
	Add(c *expr.Conjunction)
	Get(conjunctionSize int, labels expr.Assignment) []*Record
}
