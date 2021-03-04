package cnf

// PostingEntry contains the CNF ID, belongs-to flag, Disjunction ID and score
type PostingEntry struct {
	ID            int64
	DisjunctionID int8
	BelongsTo     bool
	score         int
}

// SortID returns the CNF ID
func (p *PostingEntry) SortID() int64 {
	return p.ID
}

// PostingList is a slice of unique entries
type PostingList []*PostingEntry

// PostingLists is a slice of PostingList
type PostingLists []PostingList
