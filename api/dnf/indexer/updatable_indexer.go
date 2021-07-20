package indexer

import (
	"errors"

	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
	cow "github.com/csimplestring/go-cow-loader"
)

// UpdatableIndexer
type UpdatableIndexer interface {
	MaxKSize() int
	Add(c *expr.Conjunction) error
	Delete(ID int) error
	Update(c *expr.Conjunction) error
	Get(conjunctionSize int, labels expr.Assignment) []*Record
}

// NewCopyOnWriteIndexer creats a new CopyOnWriteIndexer with given items.
// It internally uses a loader to periodically reload index.
func NewCopyOnWriteIndexer(items []*expr.Conjunction) (*CopyOnWriteIndexer, error) {
	indexV2, err := newMemIndexerV2(items)
	if err != nil {
		return nil, err
	}

	u := &CopyOnWriteIndexer{}
	u.loader = cow.New(indexV2, 300)

	return u, nil
}

// CopyOnWriteIndexer allows user to update the index. Internally it uses a loader to periodically
// reload index by applying the ADD/DEL/UPD operations.
type CopyOnWriteIndexer struct {
	loader *cow.Reloader
}

func (u *CopyOnWriteIndexer) MaxKSize() int {
	return u.loader.Reload().(*memIndexV2).MaxKSize()
}

// Add adds a new conjunction into index, this operation will be first buffered in a queue
// and be applied in a batcch way once the ticker is triggered. So the new item won't be updated immediately but with a
// lag.
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

// memIndexV2 stores all the entries in memory.
type memIndexV2 struct {
	maxKSize     int
	sizedIndexes map[int]shard
}

// NewMemIndexerV2 create a memory stored indexer.
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
	c := &memIndexV2{
		maxKSize:     k.maxKSize,
		sizedIndexes: make(map[int]shard, len(k.sizedIndexes)),
	}
	for k, v := range k.sizedIndexes {
		c.sizedIndexes[k] = v.Copy()
	}
	return c
}

func (k *memIndexV2) Apply(ops []cow.Op) error {
	for _, op := range ops {

		c := op.Context().(*expr.Conjunction)
		if op.Type() == "add" {
			k.onAdd(c)
		} else if op.Type() == "delete" {
			k.onDelete(c)
		} else if op.Type() == "update" {
			k.onUpdate(c)
		} else {
			return errors.New("unsupported op type: " + op.Type())
		}
	}
	return nil
}

func (k *memIndexV2) onAdd(c *expr.Conjunction) error {
	// validate

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

func (k *memIndexV2) onDelete(c *expr.Conjunction) error {
	return errors.New("not supported yet")
}

func (k *memIndexV2) onUpdate(c *expr.Conjunction) error {
	return errors.New("not supported yet")
}
