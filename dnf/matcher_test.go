package dnf

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

func randBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func printMemUsage() {
	bToMb := func(b uint64) uint64 {
		return b / 1024 / 1024
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func Test_kIndexTable_Add_Lot(t *testing.T) {
	k := newKIndexTable()

	for i := 0; i < 100; i++ {
		id := i

		var attrs []*Attribute
		for j := 0; j < rand.Intn(10-1)+1; j++ {
			attrs = append(attrs, &Attribute{
				Name:     randstr.String(rand.Intn(20-1) + 1),
				Values:   []string{randstr.String(rand.Intn(20-1) + 1)},
				Contains: randBool(),
			})
		}

		k.Add(NewConjunction(id, attrs))
	}

	printMemUsage()

}

func Test_kIndexTable_Add(t *testing.T) {
	k := newKIndexTable()

	k.Add(NewConjunction(
		1,
		[]*Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true},
			{Name: "state", Values: []string{"NY"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		2,
		[]*Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true},
			{Name: "gender", Values: []string{"F"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		3,
		[]*Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true},
			{Name: "gender", Values: []string{"M"}, Contains: true},
			{Name: "state", Values: []string{"CA"}, Contains: false},
		},
	))

	k.Add(NewConjunction(
		4,
		[]*Attribute{
			{Name: "state", Values: []string{"CA"}, Contains: true},
			{Name: "gender", Values: []string{"M"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		5,
		[]*Attribute{
			{Name: "age", Values: []string{"3", "4"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		6,
		[]*Attribute{
			{Name: "state", Values: []string{"CA", "NY"}, Contains: false},
		},
	))

	k.Build()

	assert.Equal(t, 2, k.MaxKSize())
	zeroIdx := k.sizedIndexes[0].(*memoryIndexer)
	assert.Equal(t, 3, len(zeroIdx.imap))
	assert.Equal(t, 6, zeroIdx.imap["state:CA"].Items[0].CID)
	assert.Equal(t, false, zeroIdx.imap["state:CA"].Items[0].Contains)
	assert.Equal(t, 6, zeroIdx.imap["state:NY"].Items[0].CID)
	assert.Equal(t, false, zeroIdx.imap["state:NY"].Items[0].Contains)
	assert.Equal(t, 6, zeroIdx.imap["null:null"].Items[0].CID)
	assert.Equal(t, true, zeroIdx.imap["null:null"].Items[0].Contains)

	oneIdx := k.sizedIndexes[1].(*memoryIndexer)
	assert.Equal(t, 2, len(oneIdx.imap))
	assert.Equal(t, 5, oneIdx.imap["age:3"].Items[0].CID)
	assert.Equal(t, true, oneIdx.imap["age:3"].Items[0].Contains)
	assert.Equal(t, 5, oneIdx.imap["age:4"].Items[0].CID)
	assert.Equal(t, true, oneIdx.imap["age:4"].Items[0].Contains)

	twoIdx := k.sizedIndexes[2].(*memoryIndexer)
	assert.Equal(t, 5, len(twoIdx.imap))
	assert.Equal(t, 1, twoIdx.imap["state:NY"].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap["state:NY"].Items[0].Contains)

	assert.Equal(t, 1, twoIdx.imap["age:3"].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap["age:3"].Items[0].Contains)
	assert.Equal(t, 2, twoIdx.imap["age:3"].Items[1].CID)
	assert.Equal(t, true, twoIdx.imap["age:3"].Items[1].Contains)
	assert.Equal(t, 3, twoIdx.imap["age:3"].Items[2].CID)
	assert.Equal(t, true, twoIdx.imap["age:3"].Items[2].Contains)

	assert.Equal(t, 2, twoIdx.imap["gender:F"].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap["gender:F"].Items[0].Contains)

	assert.Equal(t, 3, twoIdx.imap["state:CA"].Items[0].CID)
	assert.Equal(t, false, twoIdx.imap["state:CA"].Items[0].Contains)
	assert.Equal(t, 4, twoIdx.imap["state:CA"].Items[1].CID)
	assert.Equal(t, true, twoIdx.imap["state:CA"].Items[1].Contains)

	assert.Equal(t, 3, twoIdx.imap["gender:M"].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap["gender:M"].Items[0].Contains)
	assert.Equal(t, 4, twoIdx.imap["gender:M"].Items[1].CID)
	assert.Equal(t, true, twoIdx.imap["gender:M"].Items[1].Contains)

	matcher := &matcher{}
	matched := matcher.Match(k, Assignment{
		Label{Name: "age", Value: "3"},
		Label{Name: "state", Value: "CA"},
		Label{Name: "gender", Value: "M"},
	})

	assert.Equal(t, []int{4, 5}, matched)
}
