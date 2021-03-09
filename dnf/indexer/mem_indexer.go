package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/set"
)

type indexShard struct {
	invertedMap   map[uint64]postingList
	attributeMeta expr.AttributeMetadataStorer
}

func newIndexShard(attributeMeta expr.AttributeMetadataStorer) *indexShard {
	return &indexShard{
		invertedMap:   make(map[uint64]postingList),
		attributeMeta: attributeMeta,
	}
}

func (m *indexShard) toKeys(a expr.Attribute) []*key {

	keys := make([]*key, len(a.Values))
	for i, v := range a.Values {
		keys[i] = &key{
			Name:  a.Name,
			Value: v,
		}
	}
	return keys
}

func (m *indexShard) hashKey(k *key) uint64 {
	return uint64(k.Name)<<32 | uint64(k.Value)
}

func (m *indexShard) Build() error {
	for _, pList := range m.invertedMap {
		pList.sort()
	}
	return nil
}

func (m *indexShard) Get(k *key) postingList {
	h := m.hashKey(k)
	return m.invertedMap[h]
}

func (m *indexShard) createIfAbsent(hash uint64) postingList {
	v := m.get(hash)
	if v == nil {
		p := postingList{}
		m.put(hash, p)
		return p
	}
	return v
}

func (m *indexShard) get(hash uint64) postingList {
	return m.invertedMap[hash]
}

func (m *indexShard) put(hash uint64, p postingList) {
	m.invertedMap[hash] = p
}

func (m *indexShard) Add(c expr.Conjunction) error {

	for _, attr := range c.Attributes {
		for _, key := range m.toKeys(attr) {

			hash := m.hashKey(key)

			pList := m.createIfAbsent(hash)
			pList = append(pList, &postingEntry{
				CID:      c.ID,
				Contains: attr.Contains,
			})

			m.put(hash, pList)
		}
	}

	if c.GetKSize() == 0 {
		hash := m.hashKey(zKey)
		pList := m.createIfAbsent(hash)
		pList = append(pList, &postingEntry{
			CID:      c.ID,
			Contains: true,
		})
		m.put(hash, pList)
	}

	return nil
}

type memoryIndex struct {
	maxKSize      int
	attributeMeta expr.AttributeMetadataStorer
	sizedIndexes  map[int]*indexShard
}

// NewMemoryIndexer create a memory stored indexer
func NewMemoryIndexer(attributeMeta expr.AttributeMetadataStorer) Indexer {
	return &memoryIndex{
		maxKSize:      0,
		sizedIndexes:  make(map[int]*indexShard),
		attributeMeta: attributeMeta,
	}
}

func (k *memoryIndex) Add(c expr.Conjunction) {
	ksize := c.GetKSize()

	if k.maxKSize < ksize {
		k.maxKSize = ksize
	}

	kidx, exist := k.sizedIndexes[ksize]
	if !exist {
		kidx = newIndexShard(k.attributeMeta)
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

func (k *memoryIndex) getPostingLists(size int, labels expr.Assignment) []postingList {
	idx := k.sizedIndexes[size]
	if idx == nil {
		return nil
	}

	candidates := make([]postingList, 1)
	for _, label := range labels {
		name, found := k.attributeMeta.GetNameID(label.Name)
		if !found {
			return nil
		}
		value, found := k.attributeMeta.GetValueID(label.Name, label.Value)
		if !found {
			return nil
		}

		k := &key{
			Name:  name,
			Value: value,
		}
		p := idx.Get(k)
		if len(p) == 0 {
			continue
		}
		candidates = append(candidates, p)
	}
	if size == 0 {
		candidates = append(candidates, idx.Get(zKey))
	}
	return candidates
}

// Match finds the matched conjunctions given an assignment.
func (k *memoryIndex) Match(assignment expr.Assignment) []int {
	results := set.IntHashSet()

	n := min(len(assignment), k.maxKSize)

	for i := n; i >= 0; i-- {
		pLists := newPostingLists(k.getPostingLists(i, assignment))

		K := i
		if K == 0 {
			K = 1
		}
		if pLists.len() < K {
			continue
		}

		pLists.sortByCurrent()
		for pLists[K-1].current() != eolItem {
			var nextID int

			if pLists[0].current().CID == pLists[K-1].current().CID {

				if pLists[0].current().Contains == false {
					rejectID := pLists[0].current().CID
					for L := K; L <= pLists.len()-1; L++ {
						if pLists[L].current().CID == rejectID {
							pLists[L].skipTo(rejectID + 1)
						} else {
							break
						}
					}

				} else {
					results.Add(pLists[K-1].current().CID)
				}

				nextID = pLists[K-1].current().CID + 1
			} else {
				nextID = pLists[K-1].current().CID

			}

			for L := 0; L <= K-1; L++ {
				pLists[L].skipTo(nextID)
			}
			pLists.sortByCurrent()
		}

	}

	return results.ToSlice()
}
