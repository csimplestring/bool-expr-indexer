package cnf

// PostingEntry contains the CNF ID, belongs-to flag, Disjunction ID and score
type PostingEntry struct {
	ID            int64
	DisjunctionID int8
	BelongsTo     bool
	score         int
}

// PostingList is a slice of entries
type PostingList []*PostingEntry

// PostingLists is a slice of PostingList
type PostingLists []PostingList
