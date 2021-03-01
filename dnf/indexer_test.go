package dnf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortPostingList(t *testing.T) {

	p := PostingList{
		Items: []*PostingItem{
			{CID: 2, Contains: true},
			{CID: 1, Contains: true},
			{CID: 3, Contains: true},
			{CID: 1, Contains: false},
		},
	}

	p.sort()

	assert.Equal(t, int64(1), p.Items[0].CID)
	assert.Equal(t, false, p.Items[0].Contains)
	assert.Equal(t, int64(1), p.Items[1].CID)
	assert.Equal(t, true, p.Items[1].Contains)
	assert.Equal(t, int64(2), p.Items[2].CID)
	assert.Equal(t, true, p.Items[2].Contains)
	assert.Equal(t, int64(3), p.Items[3].CID)
	assert.Equal(t, true, p.Items[3].Contains)
}
