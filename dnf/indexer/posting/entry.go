package posting

import (
	"errors"
)

// EntryInt32 store conjunction-id, belongs-to flag, serving as inverted index pointing to Conjunction
type EntryInt32 interface {
	CID() uint32
	Contains() bool
	Score() uint32
}

// 1000000000
const flagMask uint32 = 1 << 31
const idMask uint32 = 0b111111111111111111111111
const scoreMask uint32 = 0b01111111000000000000000000000000
const maxID uint32 = 16777215

// EOL means end-of-list item, used as the end of a posting list.
var EOL entryInt = entryInt((1 << 31) | (0 << 24) | maxID)

// entryInt:
// the 32nd bit is the 'belongsTo' flag: 1 -> belongs, 0 -> not belongs
// the 31st ~ 25th bits are the 'score' field: [0,2^7-1], but only [0,100] is valid
// the 24th to 1st bits are the 'conjunction id' field: ranging from 0 to 16777215
type entryInt uint32

// NewEntryInt32 creates a uint32 representation of entry.
func NewEntryInt32(ID uint32, contains bool, score uint32) (EntryInt32, error) {
	if score > 100 {
		return nil, errors.New("score must be [0, 100]")
	}
	if ID > maxID {
		return nil, errors.New("ID must be [0, 16777215]")
	}

	var containsBit uint32 = 0
	if contains {
		containsBit = 1
	}

	r := (containsBit << 31) | (score << 24) | ID
	return entryInt(r), nil
}

func (e entryInt) CID() uint32 {
	return uint32(e) & idMask
}

func (e entryInt) Score() uint32 {
	return uint32(e) & scoreMask >> 24
}

func (e entryInt) Contains() bool {
	return uint32(e)&flagMask == flagMask
}
