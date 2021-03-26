package indexer

import "github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"

type Record struct {
	PostingList posting.List
	Key         string
	Value       string
}

func (r *Record) append(entry posting.EntryInt32) {
	r.PostingList = append(r.PostingList, entry)
}

func (r *Record) compact() {
	if len(r.PostingList) != cap(r.PostingList) {
		compacted := make(posting.List, len(r.PostingList))
		for i := 0; i < len(r.PostingList); i++ {
			compacted[i] = r.PostingList[i]
		}
		r.PostingList = compacted
	}
}
