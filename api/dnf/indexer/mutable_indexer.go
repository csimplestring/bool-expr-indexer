package indexer

import (
	"errors"
	"fmt"

	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
	cow "github.com/csimplestring/go-cow-loader"
)

// MutableIndexer
type MutableIndexer interface {
	MaxKSize() int
	Add(c *expr.Conjunction) error
	Delete(ID int) error
	Update(c *expr.Conjunction) error
	Get(conjunctionSize int, labels expr.Assignment) []*Record
}

var _ MutableIndexer = (*CopyOnWriteIndexer)(nil)
var _ Indexer = (*CopyOnWriteIndexer)(nil)

type internalCopyOnWriteIndexer struct {
	*MemReadOnlyIndexer
}

func (i *internalCopyOnWriteIndexer) Apply(ops []cow.Op) error {
	for _, op := range ops {
		opType := op.Type()

		switch opType {
		case "add":
			i.add(op.Context().(*expr.Conjunction))
		case "delete":
			i.delette(op.Context().(int))
		default:
			return fmt.Errorf("unsupported type: %s", opType)
		}
	}
	return errors.New("not implemented")
}

func (i *internalCopyOnWriteIndexer) add(c *expr.Conjunction) error {
	if _, exist := i.meta.forwardIdx.Get(c.ID); exist {
		return fmt.Errorf("Try to add duplicate conjunction with ID %d", c.ID)
	}

	return i.add(c)
}

func (i *internalCopyOnWriteIndexer) delette(ID int) error {
	c, exist := i.meta.forwardIdx.Get(ID)
	if !exist {
		return fmt.Errorf("Try to delete non-existing conjunction with ID %d", ID)
	}

	shard := i.sizedIndexes[c.GetKSize()]
	for _, attr := range c.Attributes {
		for _, val := range attr.Values {
			name := attr.Name
			record := shard.Get(name, val)
			record.PostingList.Remove(ID)
		}
	}

	return nil
}

func (c *internalCopyOnWriteIndexer) Copy() cow.Value {
	copiedShard := make(map[int]shard, len(c.sizedIndexes))
	for k, v := range c.sizedIndexes {
		copiedShard[k] = v.Copy()
	}

	copiedIndex := &internalCopyOnWriteIndexer{
		MemReadOnlyIndexer: &MemReadOnlyIndexer{
			maxKSize:     c.maxKSize,
			sizedIndexes: copiedShard,
		},
	}

	return copiedIndex
}

// NewCopyOnWriteIndexer creats a new CopyOnWriteIndexer with given items.
// It internally uses a loader to periodically reload index.
func NewCopyOnWriteIndexer(items []*expr.Conjunction) (*CopyOnWriteIndexer, error) {
	base, err := NewMemReadOnlyIndexer(items)
	if err != nil {
		return nil, err
	}

	idx := &internalCopyOnWriteIndexer{
		base,
	}

	u := &CopyOnWriteIndexer{}
	u.loader = cow.New(idx, 300)

	return u, nil
}

// CopyOnWriteIndexer allows user to update the index. Internally it uses a loader to periodically
// reload index by applying the ADD/DEL/UPD operations.
type CopyOnWriteIndexer struct {
	loader *cow.Reloader
}

func (u *CopyOnWriteIndexer) MaxKSize() int {
	return u.loader.Reload().(Indexer).MaxKSize()
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

func (u *CopyOnWriteIndexer) Delete(ID int) error {

	return u.loader.Accept(&IndexOp{
		OpType: "delete",
		Data:   ID,
	})
}

func (u *CopyOnWriteIndexer) Update(c *expr.Conjunction) error {

	return u.loader.Accept(&IndexOp{
		OpType: "update",
		Data:   c,
	})
}

func (u *CopyOnWriteIndexer) Get(conjunctionSize int, labels expr.Assignment) []*Record {
	idx := u.loader.Reload().(Indexer)
	return idx.Get(conjunctionSize, labels)
}
