package indexer

import (
	"hash/maphash"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
)

type indexShard struct {
	conjunctionSize int
	zeroKey         uint64
	invertedMap     map[uint64]posting.List
	hash            maphash.Hash
}

func newIndexShard(ksize int) *indexShard {
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
		invertedMap:     make(map[uint64]posting.List),
	}
}

func (m *indexShard) hashKey(name string, value string) uint64 {
	m.hash.Reset()
	m.hash.WriteString(name)
	m.hash.WriteString(value)
	return m.hash.Sum64()
}

func (m *indexShard) Build() error {

	for _, pList := range m.invertedMap {
		pList.Sort()
	}
	return nil
}

func (m *indexShard) Get(name string, value string) posting.List {
	h := m.hashKey(name, value)
	return m.invertedMap[h]
}

func (m *indexShard) createIfAbsent(hash uint64) posting.List {
	v := m.get(hash)
	if v == nil {
		p := make(posting.List, 0, 64)
		m.put(hash, p)
		return p
	}
	return v
}

func (m *indexShard) get(hash uint64) posting.List {
	return m.invertedMap[hash]
}

func (m *indexShard) put(hash uint64, p posting.List) {
	m.invertedMap[hash] = p
}

func (m *indexShard) appen(p posting.List, entry posting.EntryInt32) posting.List {
	p = append(p, entry)
	return p
}

func (m *indexShard) Add(c *expr.Conjunction) error {

	for _, attr := range c.Attributes {
		for _, value := range attr.Values {

			hash := m.hashKey(attr.Name, value)

			pList := m.createIfAbsent(hash)

			entry, err := posting.NewEntryInt32(uint32(c.ID), attr.Contains, 0)
			if err != nil {
				return err
			}

			pList = m.appen(pList, entry)

			m.put(hash, pList)
		}
	}

	if c.GetKSize() == 0 {
		pList := m.createIfAbsent(m.zeroKey)

		entry, err := posting.NewEntryInt32(uint32(c.ID), true, 0)
		if err != nil {
			return err
		}
		pList = append(pList, entry)
		m.put(m.zeroKey, pList)
	}

	return nil
}

// memoryIndex implements the Indexer interface and stores all the entries in memory.
type memoryIndex struct {
	maxKSize     int
	sizedIndexes map[int]*indexShard
}

// NewMemoryIndexer create a memory stored indexer
func NewMemoryIndexer() Indexer {
	return &memoryIndex{
		maxKSize:     0,
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
		kidx = newIndexShard(ksize)
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

func (k *memoryIndex) GetPostingLists(size int, labels expr.Assignment) []posting.List {
	idx := k.sizedIndexes[size]
	if idx == nil {
		return nil
	}

	candidates := make([]posting.List, 0, len(labels)+1)
	for _, label := range labels {

		p := idx.Get(label.Name, label.Value)
		if len(p) == 0 {
			continue
		}
		candidates = append(candidates, p)
	}
	if size == 0 {
		candidates = append(candidates, idx.Get("", ""))
	}
	return candidates
}
