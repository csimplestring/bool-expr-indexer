package indexer

type IndexOp struct {
	OpType string
	Data   interface{}
}

func (i *IndexOp) Type() string {
	return i.OpType
}

func (i *IndexOp) Context() interface{} {
	return i.Data
}
