package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortPostingList(t *testing.T) {

	p := PostingList{
		{CID: 2, Contains: true},
		{CID: 1, Contains: true},
		{CID: 3, Contains: true},
		{CID: 1, Contains: false},
	}

	p.sort()

	assert.Equal(t, 1, p[0].CID)
	assert.Equal(t, false, p[0].Contains)
	assert.Equal(t, 1, p[1].CID)
	assert.Equal(t, true, p[1].Contains)
	assert.Equal(t, 2, p[2].CID)
	assert.Equal(t, true, p[2].Contains)
	assert.Equal(t, 3, p[3].CID)
	assert.Equal(t, true, p[3].Contains)
}
