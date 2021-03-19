package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
)

// memoryIndexer implements the Indexer interface and stores all the entries in memory.
type memoryIndexer struct {
	maxKSize     int
	sizedIndexes map[int]*indexShard
}

// NewMemoryIndexer create a memory stored indexer
func NewMemoryIndexer() Indexer {
	return &memoryIndexer{
		maxKSize:     0,
		sizedIndexes: make(map[int]*indexShard),
	}
}

func (k *memoryIndexer) Add(c *expr.Conjunction) {
	ksize := c.GetKSize()

	if k.maxKSize < ksize {
		k.maxKSize = ksize
	}

	kidx, exist := k.sizedIndexes[ksize]
	if !exist {
		kidx = newIndexShard(ksize)
		k.sizedIndexes[ksize] = kidx
	}

	kidx.Add(c)
}

func (k *memoryIndexer) Build() error {
	for _, v := range k.sizedIndexes {
		if err := v.Build(); err != nil {
			return err
		}
	}
	return nil
}

func (k *memoryIndexer) MaxKSize() int {
	return k.maxKSize
}

func (k *memoryIndexer) Get(size int, labels expr.Assignment) []*Record {
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
