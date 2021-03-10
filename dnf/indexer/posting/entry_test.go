package posting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EntryInt32(t *testing.T) {

	e, err := NewEntryInt32(2, true, 100)
	assert.NoError(t, err)
	assert.Equal(t, uint32(2), e.CID())
	assert.Equal(t, true, e.Contains())
	assert.Equal(t, uint32(100), e.Score())

	e, err = NewEntryInt32(1, false, 99)
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), e.CID())
	assert.Equal(t, false, e.Contains())
	assert.Equal(t, uint32(99), e.Score())

	e, err = NewEntryInt32(1, false, 101)
	assert.Error(t, err)

	e, err = NewEntryInt32(167772150, false, 99)
	assert.Error(t, err)
}

func Benchmark_EntryInt32(b *testing.B) {

	for i := 0; i < b.N; i++ {
		e, err := NewEntryInt32(100, true, 100)
		if err == nil {
			e.CID()
			e.Contains()
			e.Score()
		}
	}
}
