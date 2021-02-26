package main

import "sort"

// KSizeIndexer shards the Indexer by conjunction size.
type KSizeIndexer interface {
	Build() error
	MaxKSize() int
	Get(size int) Indexer
}

// Indexer actually stores the reverted index: key -> posting list
type Indexer interface {
	Add(c *Conjunction) error
	Get(k *Key) *PostingList
	Build() error
}

// Key is the key representing an attribute, e.g., <age, 10>
type Key struct {
	Name  string
	Value string
	score int
}

// PostingItem store conjunction-id, belongs-to flag, serving as inverted index pointing to Conjunction
type PostingItem struct {
	CID      int64
	Contains bool
	score    int
}

// PostingList is a list of PostingItem
type PostingList struct {
	Items []*PostingItem
}

func newPostingList() *PostingList {
	return &PostingList{}
}

func (p *PostingList) append(item *PostingItem) {
	p.Items = append(p.Items, item)
}

func (p *PostingList) sort() {
	sort.Slice(p.Items[:], func(i, j int) bool {

		if p.Items[i].CID != p.Items[j].CID {
			return p.Items[i].CID < p.Items[j].CID
		}

		return !p.Items[i].Contains && p.Items[j].Contains
	})
}
