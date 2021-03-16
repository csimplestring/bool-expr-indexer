package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
)

// Indexer defines the top level indexer interface
type Indexer interface {
	Build() error
	MaxKSize() int
	Add(c *expr.Conjunction)
	GetPostingLists(conjunctionSize int, labels expr.Assignment) []posting.List
}
