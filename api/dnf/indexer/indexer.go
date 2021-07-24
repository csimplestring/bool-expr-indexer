package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
)

// Indexer defines the top level indexer interface
type Indexer interface {
	Build() error
	MaxKSize() int
	Get(conjunctionSize int, labels expr.Assignment) []*Record
}
