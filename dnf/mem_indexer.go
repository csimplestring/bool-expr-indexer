package dnf

import (
	"fmt"
	"math"
)

// zKey is zero boolean key placeholder
var zKey *Key = &Key{
	Name:  0,
	Value: 0,
	score: 0,
}

var eolItem *PostingEntry = &PostingEntry{
	score:    0,
	CID:      math.MaxInt64,
	Contains: true,
}

type memoryIndexer struct {
	imap map[string]*PostingList
}

// NewMemoryIndex creates a new memory indexer, exposed as a public api.
func NewMemoryIndex() Indexer {
	return newMemoryIndex()
}

func newMemoryIndex() *memoryIndexer {
	return &memoryIndexer{
		imap: make(map[string]*PostingList),
	}
}

func (m *memoryIndexer) createKeys(a *Attribute) []*Key {
	var keys []*Key
	for _, v := range a.Values {
		keys = append(keys, &Key{
			Name:  a.Name,
			Value: v,
		})
	}
	return keys
}

func (m *memoryIndexer) hashKey(k *Key) string {
	return fmt.Sprintf("%d:%d", k.Name, k.Value)
}

func (m *memoryIndexer) Build() error {
	for _, pList := range m.imap {
		pList.sort()
	}
	return nil
}

func (m *memoryIndexer) Get(k *Key) *PostingList {
	h := m.hashKey(k)
	return m.imap[h]
}

func (m *memoryIndexer) createIfAbsent(hash string) *PostingList {
	v := m.get(hash)
	if v == nil {
		p := newPostingList()
		m.put(hash, p)
		return p
	}
	return v
}

func (m *memoryIndexer) get(hash string) *PostingList {
	return m.imap[hash]
}

func (m *memoryIndexer) put(hash string, p *PostingList) {
	m.imap[hash] = p
}

func (m *memoryIndexer) Add(c *Conjunction) error {

	for _, attr := range c.Attributes {
		for _, key := range m.createKeys(attr) {

			hash := m.hashKey(key)

			pList := m.createIfAbsent(hash)
			pList.append(&PostingEntry{
				CID:      c.ID,
				Contains: attr.Contains,
			})

			m.put(hash, pList)
		}
	}

	if c.kSize == 0 {
		hash := m.hashKey(zKey)
		pList := m.createIfAbsent(hash)
		pList.append(&PostingEntry{
			CID:      c.ID,
			Contains: true,
		})
		m.put(hash, pList)
	}

	return nil
}

type kIndexTable struct {
	maxKSize     int
	sizedIndexes map[int]Indexer
}

func newKIndexTable() *kIndexTable {
	return &kIndexTable{
		maxKSize:     0,
		sizedIndexes: make(map[int]Indexer),
	}
}

func (k *kIndexTable) Add(c *Conjunction) {
	ksize := c.GetKSize()

	if k.maxKSize < ksize {
		k.maxKSize = ksize
	}

	kidx, exist := k.sizedIndexes[ksize]
	if !exist {
		kidx = newMemoryIndex()
		k.sizedIndexes[ksize] = kidx
	}

	kidx.Add(c)
}

func (k *kIndexTable) Build() error {
	for _, v := range k.sizedIndexes {
		if err := v.Build(); err != nil {
			return err
		}
	}
	return nil
}

func (k *kIndexTable) MaxKSize() int {
	return k.maxKSize
}

func (k *kIndexTable) GetPostingLists(size int, labels Assignment) []*PostingList {
	idx := k.sizedIndexes[size]
	if idx == nil {
		return nil
	}

	var candidates []*PostingList
	for _, label := range labels {
		k := &Key{
			Name:  label.Name,
			Value: label.Value,
		}
		p := idx.Get(k)
		if p == nil {
			continue
		}
		candidates = append(candidates, p)
	}
	if size == 0 {
		candidates = append(candidates, idx.Get(zKey))
	}
	return candidates
}
