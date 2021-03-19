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
