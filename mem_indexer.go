package main

import "sort"

// zKey is zero boolean key placeholder
var zKey *Key = &Key{
	Name:  "null",
	Value: "null",
	score: 0,
}

var eolItem *PostingItem = &PostingItem{}

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
	return k.Name + ":" + k.Value
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
			pList.append(&PostingItem{
				CID:      c.ID,
				Contains: attr.Contains,
			})

			m.put(hash, pList)
		}
	}

	if c.kSize == 0 {
		hash := m.hashKey(zKey)
		pList := m.createIfAbsent(hash)
		pList.append(&PostingItem{
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

func (k *kIndexTable) GetPostingLists(size int, labels Labels) []*PostingList {
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

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

type postingLists struct {
	c []*postingListCursor
}

func newPostingLists(l []*PostingList) *postingLists {
	var c []*postingListCursor

	for _, v := range l {
		c = append(c, newPostingListCursor(v))
	}
	return &postingLists{
		c: c,
	}
}

func (p *postingLists) len() int {
	return len(p.c)
}

func (p *postingLists) sortByCurrent() {
	sort.Slice(p.c[:], func(i, j int) bool {
		a := p.c[i].current()
		b := p.c[j].current()

		// fix
		if a == eolItem && b != eolItem {
			return false
		}
		if a != eolItem && b == eolItem {
			return true
		}

		if a.CID != b.CID {
			return a.CID < b.CID
		}

		return !a.Contains && b.Contains
	})
}

type postingListCursor struct {
	ref  *PostingList
	size int
	cur  int
}

func newPostingListCursor(ref *PostingList) *postingListCursor {
	return &postingListCursor{
		ref:  ref,
		size: len(ref.Items),
		cur:  0,
	}
}

func (p *postingListCursor) current() *PostingItem {
	if p.cur >= p.size {
		return eolItem
	}

	return p.ref.Items[p.cur]
}

func (p *postingListCursor) skipTo(ID int64) {
	i := p.cur
	for i < p.size && p.ref.Items[i].CID != ID {
		i++
	}

	p.cur = i
}

func (k *kIndexTable) Match(labels Labels) []int64 {
	var results []int64

	n := min(len(labels), k.maxKSize)

	for i := n; i >= 0; i-- {
		pLists := newPostingLists(k.GetPostingLists(i, labels))

		K := i
		if K == 0 {
			K = 1
		}
		if pLists.len() < K {
			continue
		}

		for pLists.c[K-1].current() != eolItem {
			var nextID int64

			pLists.sortByCurrent()
			if pLists.c[0].current().CID == pLists.c[K-1].current().CID {

				if pLists.c[0].current().Contains == false {
					rejectID := pLists.c[0].current().CID
					for L := K; L <= pLists.len()-1; L++ {
						if pLists.c[L].current().CID == rejectID {
							pLists.c[L].skipTo(rejectID + 1)
						} else {
							break
						}
					}

				} else {
					results = append(results, pLists.c[K-1].current().CID)
				}

				nextID = pLists.c[K-1].current().CID + 1
			} else {
				nextID = pLists.c[K-1].current().CID

			}

			for L := 0; L <= K-1; L++ {
				pLists.c[L].skipTo(nextID)
			}
			pLists.sortByCurrent()
		}

	}

	return results
}
