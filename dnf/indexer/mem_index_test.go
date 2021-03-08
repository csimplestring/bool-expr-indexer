package indexer

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
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

func Benchmark_kIndexTable_Add(b *testing.B) {
	k := NewMemoryIndexer(nil)

	for n := 0; n < 1000000; n++ {
		id := n + 1

		var attrs []expr.Attribute
		for j := 0; j < rand.Intn(10-1)+1; j++ {
			attrs = append(attrs, expr.Attribute{
				Name:     uint32(rand.Intn(200-1) + 1),
				Values:   []uint32{uint32(rand.Intn(20-1) + 1)},
				Contains: randBool(),
			})
		}

		k.Add(expr.NewConjunction(id, attrs))
	}

	printMemUsage()
}

func getTestAttribute(t *testing.T, m expr.AttributeMetadataStorer, name string, values []string, contains bool) expr.Attribute {
	a, err := m.NewAttribute(name, values, contains)
	assert.NoError(t, err)

	return a
}

func Test_kIndexTable_Match(t *testing.T) {
	metastorer := expr.NewAttributeMetadataStorer()
	metastorer.AddNameIDMapping("age", 1)
	metastorer.AddNameIDMapping("state", 2)
	metastorer.AddNameIDMapping("gender", 3)
	metastorer.AddValueIDMapping("age", "3", 3)
	metastorer.AddValueIDMapping("age", "4", 4)
	metastorer.AddValueIDMapping("state", "NY", 1)
	metastorer.AddValueIDMapping("state", "CA", 2)
	metastorer.AddValueIDMapping("gender", "F", 0)
	metastorer.AddValueIDMapping("gender", "M", 1)

	k := NewMemoryIndexer(metastorer)

	k.Add(expr.NewConjunction(
		1,
		[]expr.Attribute{
			getTestAttribute(t, metastorer, "age", []string{"3"}, true),
			getTestAttribute(t, metastorer, "state", []string{"NY"}, true),
		},
	))

	k.Add(expr.NewConjunction(
		2,
		[]expr.Attribute{
			getTestAttribute(t, metastorer, "age", []string{"3"}, true),
			getTestAttribute(t, metastorer, "gender", []string{"F"}, true),
		},
	))

	k.Add(expr.NewConjunction(
		3,
		[]expr.Attribute{
			getTestAttribute(t, metastorer, "age", []string{"3"}, true),
			getTestAttribute(t, metastorer, "gender", []string{"M"}, true),
			getTestAttribute(t, metastorer, "state", []string{"CA"}, false),
		},
	))

	k.Add(expr.NewConjunction(
		4,
		[]expr.Attribute{
			getTestAttribute(t, metastorer, "state", []string{"CA"}, true),
			getTestAttribute(t, metastorer, "gender", []string{"M"}, true),
		},
	))

	k.Add(expr.NewConjunction(
		5,
		[]expr.Attribute{
			getTestAttribute(t, metastorer, "age", []string{"3", "4"}, true),
		},
	))

	k.Add(expr.NewConjunction(
		6,
		[]expr.Attribute{
			getTestAttribute(t, metastorer, "state", []string{"CA", "NY"}, false),
		},
	))

	k.Build()

	assert.Equal(t, 2, k.MaxKSize())

	// zeroIdx := k.sizedIndexes[0].(*memoryIndexer)
	// assert.Equal(t, 3, len(zeroIdx.imap))
	// assert.Equal(t, 6, zeroIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[0].CID)
	// assert.Equal(t, false, zeroIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[0].Contains)
	// assert.Equal(t, 6, zeroIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["NY"])].Items[0].CID)
	// assert.Equal(t, false, zeroIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["NY"])].Items[0].Contains)
	// assert.Equal(t, 6, zeroIdx.imap["0:0"].Items[0].CID)
	// assert.Equal(t, true, zeroIdx.imap["0:0"].Items[0].Contains)

	// oneIdx := k.sizedIndexes[1].(*memoryIndexer)
	// assert.Equal(t, 2, len(oneIdx.imap))
	// assert.Equal(t, 5, oneIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[0].CID)
	// assert.Equal(t, true, oneIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[0].Contains)
	// assert.Equal(t, 5, oneIdx.imap[fmt.Sprintf("%d:4", attrName["age"])].Items[0].CID)
	// assert.Equal(t, true, oneIdx.imap[fmt.Sprintf("%d:4", attrName["age"])].Items[0].Contains)

	// twoIdx := k.sizedIndexes[2].(*memoryIndexer)
	// assert.Equal(t, 5, len(twoIdx.imap))
	// assert.Equal(t, 1, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["NY"])].Items[0].CID)
	// assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["NY"])].Items[0].Contains)

	// assert.Equal(t, 1, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[0].CID)
	// assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[0].Contains)
	// assert.Equal(t, 2, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[1].CID)
	// assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[1].Contains)
	// assert.Equal(t, 3, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[2].CID)
	// assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:3", attrName["age"])].Items[2].Contains)

	// assert.Equal(t, 2, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["F"])].Items[0].CID)
	// assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["F"])].Items[0].Contains)

	// assert.Equal(t, 3, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[0].CID)
	// assert.Equal(t, false, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[0].Contains)
	// assert.Equal(t, 4, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[1].CID)
	// assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["state"], stateValue["CA"])].Items[1].Contains)

	// assert.Equal(t, 3, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["M"])].Items[0].CID)
	// assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["M"])].Items[0].Contains)
	// assert.Equal(t, 4, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["M"])].Items[1].CID)
	// assert.Equal(t, true, twoIdx.imap[fmt.Sprintf("%d:%d", attrName["gender"], genderValue["M"])].Items[1].Contains)

	matched := k.Match(expr.Assignment{
		expr.Label{Name: "age", Value: "3"},
		expr.Label{Name: "state", Value: "CA"},
		expr.Label{Name: "gender", Value: "M"},
	})

	assert.ElementsMatch(t, []int{4, 5}, matched)
}
