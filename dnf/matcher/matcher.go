package matcher

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer"
	"github.com/csimplestring/bool-expr-indexer/dnf/matcher/simple"
)

type Matcher interface {
	Match(indexer indexer.Indexer, assignment expr.Assignment) []int
}

func Simple() Matcher {
	return simple.New()
}
