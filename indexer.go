package main

import (
	"sort"
)

// type ConjunctionDecomposer interface {
// 	Decompose(c *Conjunction) (*IndexKey, *)
// }

// IndexKey is the key in KIndexes table
type IndexKey struct {
	Key   string
	Value string
	score int
}

// PostingItem store conjunction-id, belongs-to flag, serving as inverted index pointing to Conjunction
type PostingItem struct {
	ConjunctionID int64
	Contains      bool
	score         int
}

// PostingList is a list of PostingItem
type PostingList []*PostingItem

type postingIndex struct {
	m map[IndexKey]PostingList
}

func newPostingIndex() *postingIndex {
	return &postingIndex{
		m: make(map[IndexKey]PostingList),
	}
}

var zeroSizeIndex IndexKey = IndexKey{
	Key:   "",
	Value: "",
	score: 0,
}

func (k *postingIndex) createKeys(a *Attribute) []*IndexKey {
	var keys []*IndexKey
	for _, v := range a.Values {
		keys = append(keys, &IndexKey{
			Key:   a.Key,
			Value: v,
		})
	}
	return keys
}

func (k *postingIndex) Add(c *Conjunction) {

	for _, attr := range c.Attributes {
		for _, key := range k.createKeys(attr) {
			pList := k.m[*key]
			pList = append(pList, &PostingItem{
				ConjunctionID: c.ID,
				Contains:      attr.Contains,
			})
			k.m[*key] = pList
		}
	}

	if c.kSize == 0 {
		pList := k.m[zeroSizeIndex]
		pList = append(pList, &PostingItem{
			ConjunctionID: c.ID,
			Contains:      true,
		})
		k.m[zeroSizeIndex] = pList
	}
}

type kIndexTable struct {
	maxKSize int
	store    map[int]*postingIndex
}

func newKIndexTable() *kIndexTable {
	return &kIndexTable{
		maxKSize: 0,
		store:    make(map[int]*postingIndex),
	}
}

func (k *kIndexTable) Add(c *Conjunction) {
	ksize := c.GetKSize()

	if k.maxKSize < ksize {
		k.maxKSize = ksize
	}

	kidx, exist := k.store[ksize]
	if !exist {
		kidx = newPostingIndex()
		k.store[ksize] = kidx
	}

	kidx.Add(c)
}

func sortPostingList(p PostingList) {
	sort.Slice(p[:], func(i, j int) bool {

		if p[i].ConjunctionID != p[j].ConjunctionID {
			return p[i].ConjunctionID < p[j].ConjunctionID
		}

		return !p[i].Contains && p[j].Contains
	})
}

func (k *kIndexTable) Build() {
	for _, v := range k.store {
		for _, pList := range v.m {
			sortPostingList(pList)
		}
	}
}

func (k *kIndexTable) MaxKSize() int {
	return k.maxKSize
}
