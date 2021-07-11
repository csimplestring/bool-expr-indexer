package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
)

// Indexer defines the top level indexer interface
type Indexer interface {
	Build() error
	MaxKSize() int
	Add(c *expr.Conjunction) error
	Get(conjunctionSize int, labels expr.Assignment) []*Record
}

// memIndexer implements the Indexer interface and stores all the entries in memory.
type memIndexer struct {
	maxKSize     int
	sizedIndexes map[int]shard
}

// NewMemIndexer create a memory stored indexer. This kind of indexer is thread-safe.
func NewMemIndexer() Indexer {
	return &memIndexer{
		maxKSize:     0,
		sizedIndexes: make(map[int]shard),
	}
}

// Add conjunction into indexer.
func (k *memIndexer) Add(c *expr.Conjunction) error {
	ksize := c.GetKSize()

	if k.maxKSize < ksize {
		k.maxKSize = ksize
	}

	kidx, exist := k.sizedIndexes[ksize]
	if !exist {
		kidx = newMapShard(ksize)
		k.sizedIndexes[ksize] = kidx
	}

	return kidx.Add(c)
}

// Build finalise the build-up of indexer by calling the Build on each shards.
func (k *memIndexer) Build() error {
	for _, v := range k.sizedIndexes {
		if err := v.Build(); err != nil {
			return err
		}
	}
	return nil
}

// MaxKSize returns the max K-size of the conjunctions stored in indexer.
func (k *memIndexer) MaxKSize() int {
	return k.maxKSize
}

// Get returns the list of Record, based on size and labels.
func (k *memIndexer) Get(size int, labels expr.Assignment) []*Record {
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
