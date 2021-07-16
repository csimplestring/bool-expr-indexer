package indexer

import "github.com/csimplestring/bool-expr-indexer/api/dnf/expr"

type IndexOp struct {
	OpType string
	Data   *expr.Conjunction
}

func (i *IndexOp) Type() string {
	return i.OpType
}

func (i *IndexOp) Context() interface{} {
	return i.Data
}
