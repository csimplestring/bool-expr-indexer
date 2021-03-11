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

// key is the key representing an attribute, e.g., <age, 10>
type key struct {
	Name  string
	Value string
	score int
}

// zKey is zero key placeholder
var zKey *key = &key{
	Name:  "",
	Value: "",
	score: 0,
}
