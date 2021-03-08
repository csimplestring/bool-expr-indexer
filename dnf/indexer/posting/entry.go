package posting

import "math"

// EOL means end-of-list item, used as the end of a posting list.
var EOL *Entry = &Entry{
	score:    0,
	CID:      math.MaxInt64,
	Contains: true,
}

// Entry store conjunction-id, belongs-to flag, serving as inverted index pointing to Conjunction
type Entry struct {
	CID      int
	Contains bool
	score    int
}
