package posting

import (
	"errors"
)

const flagMask uint32 = 1 << 31
const idMask uint32 = 0b111111111111111111111111
const scoreMask uint32 = 0b01111111000000000000000000000000
const maxID uint32 = 16777215

var EOL EntryInt32 = EntryInt32((1 << 31) | (0 << 24) | maxID)

// EOL means end-of-list item, used as the end of a posting list.
// func EOL() EntryInt32 {
// 	return EntryInt32((1 << 31) | (0 << 24) | maxID)
// }

// EntryInt32 is
// the 32nd bit is the 'belongsTo' flag: 1 -> belongs, 0 -> not belongs
// the 31st ~ 25th bits are the 'score' field: [0,2^7-1], but only [0,100] is valid
// the 24th to 1st bits are the 'conjunction id' field: ranging from 0 to 16777215
type EntryInt32 uint32

// NewEntryInt32 creates a uint32 representation of entry.
func NewEntryInt32(ID uint32, contains bool, score uint32) (EntryInt32, error) {
	if score > 100 {
		return 0, errors.New("score must be [0, 100]")
	}
	if ID > maxID {
		return 0, errors.New("ID must be [0, 16777215]")
	}

	var containsBit uint32 = 0
	if contains {
		containsBit = 1
	}

	r := (containsBit << 31) | (score << 24) | ID
	return EntryInt32(r), nil
}

// CID returns the conjunction id
func (e EntryInt32) CID() uint32 {
	return uint32(e) & idMask
}

// Score returns the score
func (e EntryInt32) Score() uint32 {
	return uint32(e) & scoreMask >> 24
}

// Contains returns the belongs-to flag
func (e EntryInt32) Contains() bool {
	return uint32(e)&flagMask == flagMask
}
