package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
)

// KSizeIndexer shards the Indexer by conjunction size.
type KSizeIndexer interface {
	Build() error
	MaxKSize() int
	Add(c *expr.Conjunction)
	Match(assignment expr.Assignment) []int
}

// Indexer actually stores the reverted index: key -> posting list
type Indexer interface {
	Add(c *expr.Conjunction) error
	Get(k *key) *PostingList
	Build() error
}

// key is the key representing an attribute, e.g., <age, 10>
type key struct {
	Name  uint32
	Value uint32
	score int
}

// zKey is zero key placeholder
var zKey *key = &key{
	Name:  0,
	Value: 0,
	score: 0,
}
