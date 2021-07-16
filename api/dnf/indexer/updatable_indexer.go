package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
	cow "github.com/csimplestring/go-cow-loader"
)

type UpdatableIndexer interface {
	MaxKSize() int
	Add(c *expr.Conjunction) error
	Delete(c *expr.Conjunction) error
	Update(c *expr.Conjunction) error
	Get(conjunctionSize int, labels expr.Assignment) []*Record
}

func NewCopyOnWriteIndexer(items []*expr.Conjunction) (*CopyOnWriteIndexer, error) {
	indexV2, err := newMemIndexerV2(items)
	if err != nil {
		return nil, err
	}

	u := &CopyOnWriteIndexer{}
	u.loader = cow.New(indexV2, 300)
	go u.loader.Start()

	return u, nil
}

type CopyOnWriteIndexer struct {
	loader *cow.Reloader
}

func (u *CopyOnWriteIndexer) MaxKSize() int {
	return u.loader.Reload().(*memIndexV2).MaxKSize()
}

func (u *CopyOnWriteIndexer) Add(c *expr.Conjunction) error {
	return u.loader.Accept(&IndexOp{
		OpType: "add",
		Data:   c,
	})
}

func (u *CopyOnWriteIndexer) Delete(c *expr.Conjunction) error {
	return u.loader.Accept(&IndexOp{
		OpType: "delete",
		Data:   c,
	})
}
func (u *CopyOnWriteIndexer) Update(c *expr.Conjunction) error {
	return u.loader.Accept(&IndexOp{
		OpType: "update",
		Data:   c,
	})
}
func (u *CopyOnWriteIndexer) Get(conjunctionSize int, labels expr.Assignment) []*Record {
	idx := u.loader.Reload().(*memIndexV2)
	return idx.Get(conjunctionSize, labels)
}

// memIndexV2 implements the Indexer interface and stores all the entries in memory.
type memIndexV2 struct {
	maxKSize     int
	sizedIndexes map[int]shard
}

// NewMemIndexerV2 create a memory stored indexer. This kind of indexer is thread-safe.
func newMemIndexerV2(items []*expr.Conjunction) (*memIndexV2, error) {

	m := &memIndexV2{
		maxKSize:     0,
		sizedIndexes: make(map[int]shard),
	}

	for _, c := range items {
		ksize := c.GetKSize()

		if m.maxKSize < ksize {
			m.maxKSize = ksize
		}

		kidx, exist := m.sizedIndexes[ksize]
		if !exist {
			kidx = newMapShard(ksize)
			m.sizedIndexes[ksize] = kidx
		}

		if err := kidx.Add(c); err != nil {
			return nil, err
		}
	}

	for _, v := range m.sizedIndexes {
		if err := v.Build(); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// MaxKSize returns the max K-size of the conjunctions stored in indexer.
func (k *memIndexV2) MaxKSize() int {
	return k.maxKSize
}

// Get returns the list of Record, based on size and labels.
func (k *memIndexV2) Get(size int, labels expr.Assignment) []*Record {
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

func (k *memIndexV2) Copy() cow.Value {

	return nil
}

func (k *memIndexV2) Apply(ops []cow.Op) error {
	return nil
}
