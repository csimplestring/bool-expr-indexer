package main

import (
	"fmt"
	"sort"
)

// Attribute is the pair of key-value, e.g., age:10, representing 'belongs to'
// The value here is discrete, but for range values, like age < 40, we can convert it into multiple pairs:
// age < 40 ---> age = 10, age = 20, age = 30, the granularity is 10.
// TODO: for longer range values such as dates, using a hierarchy structure to convert
type Attribute struct {
	Key       string
	Value     []string
	BelongsTo bool
}

// Conjunction consists of a slice of Attributes, which are combined with 'AND' logic.
type Conjunction struct {
	ID         int64
	Attributes []*Attribute
	kSize      int
}

// GetKSize ...
func (c *Conjunction) GetKSize() int {
	return c.kSize
}

// NewConjunction creates a new Conjunction
func NewConjunction(ID int64, attrs []*Attribute) *Conjunction {
	ksize := 0
	for _, a := range attrs {
		if a.BelongsTo {
			ksize++
		}
	}
	return &Conjunction{
		ID:         ID,
		Attributes: attrs,
		kSize:      ksize,
	}
}

// IndexKey is the key in KIndexes table
type IndexKey struct {
	Key   string
	Value string
	score int
}

// PostingItem store conjunction-id, belongs-to flag, serving as inverted index pointing to Conjunction
type PostingItem struct {
	ConjunctionID int64
	BelongsTo     bool
	score         int
}

// PostingList is a list of PostingItem
type PostingList []PostingItem

type kIndex struct {
	m map[string]PostingList
}

func newKIndex() *kIndex {
	return &kIndex{
		m: make(map[string]PostingList),
	}
}

var zeroSizeIndex string = "null:null:0"

func (k *kIndex) createKeys(a *Attribute) []string {
	var keys []string
	for _, v := range a.Value {
		keys = append(keys, fmt.Sprintf("%s,%s", a.Key, v))
	}
	return keys
}

func (k *kIndex) Add(c *Conjunction) {

	for _, attr := range c.Attributes {
		for _, key := range k.createKeys(attr) {
			pList := k.m[key]
			pList = append(pList, PostingItem{
				ConjunctionID: c.ID,
				BelongsTo:     attr.BelongsTo,
			})
			k.m[key] = pList
		}
	}

	if c.kSize == 0 {
		pList := k.m[zeroSizeIndex]
		pList = append(pList, PostingItem{
			ConjunctionID: c.ID,
			BelongsTo:     true,
		})
		k.m[zeroSizeIndex] = pList
	}
}

type kIndexTable struct {
	maxKSize int
	store    map[int]*kIndex
}

func newKIndexTable() *kIndexTable {
	return &kIndexTable{
		maxKSize: 0,
		store:    make(map[int]*kIndex),
	}
}

func (k *kIndexTable) Add(c *Conjunction) {
	ksize := c.GetKSize()

	if k.maxKSize < ksize {
		k.maxKSize = ksize
	}

	kidx, exist := k.store[ksize]
	if !exist {
		kidx = newKIndex()
		k.store[ksize] = kidx
	}

	kidx.Add(c)
}

func (k *kIndexTable) Build() {
	for _, v := range k.store {
		for _, pList := range v.m {
			sort.Slice(pList[:], func(i, j int) bool {
				if pList[i].ConjunctionID == pList[j].ConjunctionID {
					return !pList[i].BelongsTo && pList[j].BelongsTo
				}

				return pList[i].ConjunctionID < pList[j].ConjunctionID
			})
		}
	}
}

func (k *kIndexTable) MaxKSize() int {
	return k.maxKSize
}
