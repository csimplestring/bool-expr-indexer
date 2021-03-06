package indexer

import (
	"fmt"

	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
	cmap "github.com/orcaman/concurrent-map"
)

var _ Indexer = (*MemReadOnlyIndexer)(nil)

// MemReadOnlyIndexer implements the Indexer interface and stores all the entries in memory.
type MemReadOnlyIndexer struct {
	meta         *metadata
	maxKSize     int
	sizedIndexes map[int]shard
}

// NewMemReadOnlyIndexer create a memory stored indexer. This kind of indexer is thread-safe.
func NewMemReadOnlyIndexer(items []*expr.Conjunction) (*MemReadOnlyIndexer, error) {

	m := &MemReadOnlyIndexer{
		meta: &metadata{
			forwardIdx: &forwardIndex{
				cmap.New(),
			},
		},
		maxKSize:     0,
		sizedIndexes: make(map[int]shard),
	}

	for _, item := range items {
		if err := m.add(item); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// Add conjunction into indexer.
func (k *MemReadOnlyIndexer) add(c *expr.Conjunction) error {
	if _, exist := k.meta.forwardIdx.Get(c.ID); exist {
		return fmt.Errorf("duplicate conjunction with ID: %d", c.ID)
	}

	ksize := c.GetKSize()

	if k.maxKSize < ksize {
		k.maxKSize = ksize
	}

	kidx, exist := k.sizedIndexes[ksize]
	if !exist {
		kidx = newMapShard(ksize)
		k.sizedIndexes[ksize] = kidx
	}

	if err := kidx.Add(c); err != nil {
		return err
	}
	k.meta.forwardIdx.Set(c.ID, c)
	return nil
}

// Build finalise the build-up of indexer by calling the Build on each shards.
func (k *MemReadOnlyIndexer) Build() error {
	for _, v := range k.sizedIndexes {
		if err := v.Build(); err != nil {
			return err
		}
	}
	return nil
}

// MaxKSize returns the max K-size of the conjunctions stored in indexer.
func (k *MemReadOnlyIndexer) MaxKSize() int {
	return k.maxKSize
}

// Get returns the list of Record, based on size and labels.
func (k *MemReadOnlyIndexer) Get(size int, labels expr.Assignment) []*Record {
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
