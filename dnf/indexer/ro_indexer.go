package indexer

import (
	"errors"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
)

// readOnlyMemIndexer implements the Indexer interface and stores all the entries in memory.
type readOnlyMemIndexer struct {
	isDone       bool
	maxKSize     int
	sizedIndexes map[int]*indexShard
}

// NewReadOnlyMemIndexer create a memory stored indexer. This kind of indexer is not thread-safe but give the best performance!
func NewReadOnlyMemIndexer() Indexer {
	return &readOnlyMemIndexer{
		isDone:       false,
		maxKSize:     0,
		sizedIndexes: make(map[int]*indexShard),
	}
}

func (k *readOnlyMemIndexer) Add(c *expr.Conjunction) error {
	if k.isDone {
		return errors.New("Write operation is not supported in read-only indexer")
	}
	ksize := c.GetKSize()

	if k.maxKSize < ksize {
		k.maxKSize = ksize
	}

	kidx, exist := k.sizedIndexes[ksize]
	if !exist {
		kidx = newIndexShard(ksize)
		k.sizedIndexes[ksize] = kidx
	}

	return kidx.Add(c)
}

func (k *readOnlyMemIndexer) Build() error {
	for _, v := range k.sizedIndexes {
		if err := v.Build(); err != nil {
			return err
		}
	}
	k.isDone = true
	return nil
}

func (k *readOnlyMemIndexer) MaxKSize() int {
	return k.maxKSize
}

func (k *readOnlyMemIndexer) Get(size int, labels expr.Assignment) []*Record {
	idx := k.sizedIndexes[size]
	if idx == nil {
		return nil
	}

	candidates := make([]*Record, 0, len(labels)+1)
	for _, label := range labels {

		p := idx.Get(label.Name, label.Value)
		if p == nil {
			continue
		}
		candidates = append(candidates, p)
	}
	if size == 0 {
		candidates = append(candidates, idx.Get("", ""))
	}
	return candidates
}
