package indexer

import (
	"testing"

	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
	"github.com/stretchr/testify/assert"
)

func Test_postingList(t *testing.T) {

	e1, _ := posting.NewEntryInt32(3, true, 0)
	e2, _ := posting.NewEntryInt32(1, true, 0)
	e3, _ := posting.NewEntryInt32(2, true, 0)
	e4, _ := posting.NewEntryInt32(2, false, 0)

	p := postingList{e1, e2, e3, e4}

	p.sort()

	assert.Equal(t, uint32(1), p[0].CID())
	assert.Equal(t, uint32(2), p[1].CID())
	assert.Equal(t, uint32(2), p[2].CID())
	assert.Equal(t, uint32(3), p[3].CID())

	assert.Equal(t, true, p[0].Contains())
	assert.Equal(t, false, p[1].Contains())
	assert.Equal(t, true, p[2].Contains())
	assert.Equal(t, true, p[3].Contains())
}

func Test_plistIter(t *testing.T) {
	e1, _ := posting.NewEntryInt32(1, true, 0)
	e2, _ := posting.NewEntryInt32(2, true, 0)
	e3, _ := posting.NewEntryInt32(3, true, 0)
	e4, _ := posting.NewEntryInt32(3, false, 0)

	p := postingList{e1, e2, e3, e4}

	iter := newIterator(p)
	assert.Equal(t, e1, iter.current())

	iter.skipTo(2)
	assert.Equal(t, e2, iter.current())

	iter.skipTo(3)
	assert.Equal(t, e3, iter.current())

	iter.skipTo(3)
	assert.Equal(t, e3, iter.current())

	iter.skipTo(4)
	assert.Equal(t, posting.EOL, iter.current())
}

func Test_postingLists(t *testing.T) {

	e1, _ := posting.NewEntryInt32(1, true, 0)
	e2, _ := posting.NewEntryInt32(2, true, 0)
	e3, _ := posting.NewEntryInt32(3, true, 0)
	e4, _ := posting.NewEntryInt32(3, false, 0)

	p1 := postingList{e1}
	p2 := postingList{e2}
	p3 := postingList{e3}
	p4 := postingList{e4}

	plists := newPostingLists([]postingList{p4, p3, p2, p1})
	plists.sortByCurrent()

	assert.Equal(t, p1, plists[0].ref)
	assert.Equal(t, p2, plists[1].ref)
	assert.Equal(t, p4, plists[2].ref)
	assert.Equal(t, p3, plists[3].ref)
}
