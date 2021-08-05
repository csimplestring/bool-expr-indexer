package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
)

// Indexer defines the top level indexer interface
type Indexer interface {
	MaxKSize() int
	Get(conjunctionSize int, labels expr.Assignment) []*Record
}
