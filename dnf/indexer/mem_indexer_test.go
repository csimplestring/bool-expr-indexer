package indexer

import (
	"testing"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/stretchr/testify/assert"
)

func Test_indexShard_toKeys(t *testing.T) {
	m := &indexShard{}

	k1 := m.toKeys(expr.Attribute{Name: 1, Values: []uint32{1}})
	assert.Equal(t, []*key{{Name: 1, Value: 1}}, k1)

	k2 := m.toKeys(expr.Attribute{Name: 1, Values: []uint32{1, 2}})
	assert.ElementsMatch(t, []*key{{Name: 1, Value: 1}, {Name: 1, Value: 2}}, k2)
}
