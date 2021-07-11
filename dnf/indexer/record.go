package indexer

import "github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"

// Record represents a key-value entry in indexer shard.
type Record struct {
	PostingList posting.List
	Key         string
	Value       string
}

// append appends entry to r.
func (r *Record) append(entry posting.EntryInt32) {
	r.PostingList = append(r.PostingList, entry)
}

// compact shrinks the posting list to avoid empty slot in slice.
func (r *Record) compact() {
	if len(r.PostingList) != cap(r.PostingList) {
		compacted := make(posting.List, len(r.PostingList))
		for i := 0; i < len(r.PostingList); i++ {
			compacted[i] = r.PostingList[i]
		}
		r.PostingList = compacted
	}
}
