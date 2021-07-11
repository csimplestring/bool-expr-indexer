package posting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_postingList(t *testing.T) {

	e1, _ := NewEntryInt32(3, true, 0)
	e2, _ := NewEntryInt32(1, true, 0)
	e3, _ := NewEntryInt32(2, true, 0)
	e4, _ := NewEntryInt32(2, false, 0)

	p := List{e1, e2, e3, e4}

	p.Sort()

	assert.Equal(t, uint32(1), p[0].CID())
	assert.Equal(t, uint32(2), p[1].CID())
	assert.Equal(t, uint32(2), p[2].CID())
	assert.Equal(t, uint32(3), p[3].CID())

	assert.Equal(t, true, p[0].Contains())
	assert.Equal(t, false, p[1].Contains())
	assert.Equal(t, true, p[2].Contains())
	assert.Equal(t, true, p[3].Contains())
}
