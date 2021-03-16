package simple

import (
	"testing"

	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
	"github.com/stretchr/testify/assert"
)

func Test_plistIter(t *testing.T) {
	e1, _ := posting.NewEntryInt32(1, true, 0)
	e2, _ := posting.NewEntryInt32(2, true, 0)
	e3, _ := posting.NewEntryInt32(3, true, 0)
	e4, _ := posting.NewEntryInt32(3, false, 0)

	p := posting.List{e1, e2, e3, e4}

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

	p1 := posting.List{e1}
	p2 := posting.List{e2}
	p3 := posting.List{e3}
	p4 := posting.List{e4}

	plists := newPostingLists([]posting.List{p4, p3, p2, p1})
	plists.sortByCurrent()

	assert.Equal(t, p1, plists[0].ref)
	assert.Equal(t, p2, plists[1].ref)
	assert.Equal(t, p4, plists[2].ref)
	assert.Equal(t, p3, plists[3].ref)
}
