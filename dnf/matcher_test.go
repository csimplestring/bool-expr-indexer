package dnf

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

// func Test_kIndexTable_Add_Lot(t *testing.T) {
// 	k := newKIndexTable()

// 	for i := 0; i < 1000000; i++ {
// 		id := i

// 		var attrs []*Attribute
// 		for j := 0; j < rand.Intn(10-1)+1; j++ {
// 			attrs = append(attrs, &Attribute{
// 				Name:     rand.Intn(200-1) + 1,
// 				Values:   []int{rand.Intn(20-1) + 1},
// 				Contains: randBool(),
// 			})
// 		}

// 		k.Add(NewConjunction(id, attrs))
// 	}

// 	printMemUsage()
// }

func Test_kIndexTable_Match(t *testing.T) {
	k := newKIndexTable()

	attrName := map[string]int{
		"age": 1, "state": 2, "gender": 3,
	}
	stateValue := map[string]int{
		"NY": 1, "CA": 2,
	}
	genderValue := map[string]int{
		"F": 0, "M": 1,
	}

	k.Add(NewConjunction(
		1,
		[]*Attribute{
			{Name: attrName["age"], Values: []int{3}, Contains: true},
			{Name: attrName["state"], Values: []int{stateValue["NY"]}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		2,
		[]*Attribute{
			{Name: attrName["age"], Values: []int{3}, Contains: true},
			{Name: attrName["gender"], Values: []int{genderValue["F"]}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		3,
		[]*Attribute{
			{Name: attrName["age"], Values: []int{3}, Contains: true},
			{Name: attrName["gender"], Values: []int{genderValue["M"]}, Contains: true},
			{Name: attrName["state"], Values: []int{stateValue["CA"]}, Contains: false},
		},
	))

	k.Add(NewConjunction(
		4,
		[]*Attribute{
			{Name: attrName["state"], Values: []int{stateValue["CA"]}, Contains: true},
			{Name: attrName["gender"], Values: []int{genderValue["M"]}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		5,
		[]*Attribute{
			{Name: attrName["age"], Values: []int{3, 4}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		6,
		[]*Attribute{
			{Name: attrName["state"], Values: []int{stateValue["CA"], stateValue["NY"]}, Contains: false},
		},
	))

	k.Build()

	assert.Equal(t, 2, k.MaxKSize())
	zeroIdx := k.sizedIndexes[0].(*memoryIndexer)
	assert.Equal(t, 3, len(zeroIdx.imap))
	assert.Equal(t, 6, zeroIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[0].CID)
	assert.Equal(t, false, zeroIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[0].Contains)
	assert.Equal(t, 6, zeroIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["NY"])].Items[0].CID)
	assert.Equal(t, false, zeroIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["NY"])].Items[0].Contains)
	assert.Equal(t, 6, zeroIdx.imap["0:0"].Items[0].CID)
	assert.Equal(t, true, zeroIdx.imap["0:0"].Items[0].Contains)

	oneIdx := k.sizedIndexes[1].(*memoryIndexer)
	assert.Equal(t, 2, len(oneIdx.imap))
	assert.Equal(t, 5, oneIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[0].CID)
	assert.Equal(t, true, oneIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[0].Contains)
	assert.Equal(t, 5, oneIdx.imap[fmt.Sprintf("%d:4", attrName["age"])].Items[0].CID)
	assert.Equal(t, true, oneIdx.imap[fmt.Sprintf("%d:4", attrName["age"])].Items[0].Contains)

	twoIdx := k.sizedIndexes[2].(*memoryIndexer)
	assert.Equal(t, 5, len(twoIdx.imap))
	assert.Equal(t, 1, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["NY"])].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["NY"])].Items[0].Contains)

	assert.Equal(t, 1, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[0].Contains)
	assert.Equal(t, 2, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[1].CID)
	assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[1].Contains)
	assert.Equal(t, 3, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[2].CID)
	assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[2].Contains)

	assert.Equal(t, 2, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["F"])].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["F"])].Items[0].Contains)

	assert.Equal(t, 3, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[0].CID)
	assert.Equal(t, false, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[0].Contains)
	assert.Equal(t, 4, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[1].CID)
	assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[1].Contains)

	assert.Equal(t, 3, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["M"])].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["M"])].Items[0].Contains)
	assert.Equal(t, 4, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["M"])].Items[1].CID)
	assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["M"])].Items[1].Contains)

	matcher := &matcher{}
	matched := matcher.Match(k, Assignment{
		Label{Name: attrName["age"], Value: 3},
		Label{Name: attrName["state"], Value: stateValue["CA"]},
		Label{Name: attrName["gender"], Value: genderValue["M"]},
	})

	assert.Equal(t, []int{4, 5}, matched)
}
