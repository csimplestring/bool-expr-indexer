package indexer

import (
	"hash/maphash"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
)

type indexShard struct {
	conjunctionSize int
	zeroKey         uint64
	invertedMap     map[uint64]*Record
	scorer          Scorer
	hash            maphash.Hash
}

func newIndexShard(ksize int, scorer Scorer) *indexShard {
	// init a hasher
	var hasher maphash.Hash
	hasher.Reset()
	hasher.WriteString("")
	hasher.WriteString("")
	zeroKey := hasher.Sum64()

	return &indexShard{
		zeroKey:         zeroKey,
		conjunctionSize: ksize,
		hash:            hasher,
		scorer:          scorer,
		invertedMap:     make(map[uint64]*Record),
	}
}

func (m *indexShard) hashKey(name string, value string) uint64 {
	m.hash.Reset()
	m.hash.WriteString(name)
	m.hash.WriteString(value)
	return m.hash.Sum64()
}

func (m *indexShard) Build() error {

	for _, r := range m.invertedMap {
		r.PostingList.Sort()
	}
	return nil
}

func (m *indexShard) Get(name string, value string) *Record {
	h := m.hashKey(name, value)
	return m.invertedMap[h]
}

func (m *indexShard) createIfAbsent(hash uint64, name, value string) *Record {
	v := m.invertedMap[hash]

	if v == nil {
		r := &Record{
			PostingList: make(posting.List, 0, 64),
			Key:         name,
			Value:       value,
		}
		m.invertedMap[hash] = r
		return r
	}
	return v
}

// func (m *indexShard) get(hash uint64) *Record {
// 	return m.invertedMap[hash]
// }

// func (m *indexShard) appen(p posting.List, entry posting.EntryInt32) posting.List {
// 	p = append(p, entry)
// 	return p
// }

func (m *indexShard) Add(c *expr.Conjunction) error {

	for _, attr := range c.Attributes {
		for _, value := range attr.Values {

			hash := m.hashKey(attr.Name, value)

			r := m.createIfAbsent(hash, attr.Name, value)

			entry, err := posting.NewEntryInt32(uint32(c.ID), attr.Contains, 0)
			if err != nil {
				return err
			}

			r.append(entry)

			m.invertedMap[hash] = r
		}
	}

	if c.GetKSize() == 0 {
		r := m.createIfAbsent(m.zeroKey, "", "")

		entry, err := posting.NewEntryInt32(uint32(c.ID), true, 0)
		if err != nil {
			return err
		}
		r.append(entry)
		m.invertedMap[m.zeroKey] = r
	}

	return nil
}

// memoryIndex implements the Indexer interface and stores all the entries in memory.
type memoryIndex struct {
	maxKSize     int
	scorer       Scorer
	sizedIndexes map[int]*indexShard
}

// NewMemoryIndexer create a memory stored indexer
func NewMemoryIndexer(scorer Scorer) Indexer {
	return &memoryIndex{
		maxKSize:     0,
		scorer:       scorer,
		sizedIndexes: make(map[int]*indexShard),
	}
}

func (k *memoryIndex) Add(c *expr.Conjunction) {
	ksize := c.GetKSize()

	if k.maxKSize < ksize {
		k.maxKSize = ksize
	}

	kidx, exist := k.sizedIndexes[ksize]
	if !exist {
		kidx = newIndexShard(ksize, k.scorer)
		k.sizedIndexes[ksize] = kidx
	}

	kidx.Add(c)
}

func (k *memoryIndex) Build() error {
	for _, v := range k.sizedIndexes {
		if err := v.Build(); err != nil {
			return err
		}
	}
	return nil
}

func (k *memoryIndex) MaxKSize() int {
	return k.maxKSize
}

func (k *memoryIndex) Get(size int, labels expr.Assignment) []*Record {
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
