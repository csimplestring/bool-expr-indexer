package indexer

import (
	"hash/maphash"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
	"github.com/csimplestring/bool-expr-indexer/set"
)

type indexShard struct {
	conjunctionSize int
	zeroKey         uint64
	invertedMap     map[uint64]postingList
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
		invertedMap:     make(map[uint64]postingList),
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
		pList.sort()
	}
	return nil
}

func (m *indexShard) Get(name string, value string) postingList {
	h := m.hashKey(name, value)
	return m.invertedMap[h]
}

func (m *indexShard) createIfAbsent(hash uint64) postingList {
	v := m.get(hash)
	if v == nil {
		p := make(postingList, 0, 64)
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

func (m *indexShard) appen(p postingList, entry posting.EntryInt32) postingList {
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

func (k *memoryIndex) getPostingLists(size int, labels expr.Assignment) []postingList {
	idx := k.sizedIndexes[size]
	if idx == nil {
		return nil
	}

	candidates := make([]postingList, 0, len(labels)+1)
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
		for pLists[K-1].current() != posting.EOL() {
			var nextID uint32

			if pLists[0].current().CID() == pLists[K-1].current().CID() {

				if pLists[0].current().Contains() == false {
					rejectID := pLists[0].current().CID()
					for L := K; L <= pLists.len()-1; L++ {
						if pLists[L].current().CID() == rejectID {
							pLists[L].skipTo(int(rejectID) + 1)
						} else {
							break
						}
					}

				} else {
					results.Add(int(pLists[K-1].current().CID()))
				}

				nextID = pLists[K-1].current().CID() + 1
			} else {
				nextID = pLists[K-1].current().CID()

			}

			for L := 0; L <= K-1; L++ {
				pLists[L].skipTo(int(nextID))
			}
			pLists.sortByCurrent()
		}

	}

	return results.ToSlice()
}
