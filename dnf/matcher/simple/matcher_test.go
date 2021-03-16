package simple

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer"
	"github.com/dchest/uniuri"
	"github.com/stretchr/testify/assert"
)

func randBool() bool {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	if n <= 1 {
		return false
	}
	return true
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

func getTestAttributes() (map[string][]string, []string) {

	m := make(map[string][]string)

	names := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		nameLen := rand.Intn(8) + 1
		name := uniuri.NewLen(nameLen)

		valuesLen := rand.Intn(20) + 1
		values := make([]string, valuesLen)
		for j := 0; j < valuesLen; j++ {
			values[j] = uniuri.NewLen(valuesLen)
		}

		m[name] = values
		names[i] = name
	}

	return m, names
}

func getTestAssignment(n int, attrs map[string][]string, names []string) expr.Assignment {
	labels := make([]expr.Label, n)
	for i := 0; i < n; i++ {
		name := names[rand.Intn(len(names))]
		values := attrs[name]

		value := values[0]
		flag := randBool()
		if flag {
			value = ""
		}

		labels[i] = expr.Label{
			Name:  name,
			Value: value,
		}
	}
	return labels
}

func getIndexerAndAssignment(conjunctionNum, assignmentNum, assignmentAvgSize int) (indexer.Indexer, []expr.Assignment) {
	testAttrs, testAttrNames := getTestAttributes()

	k := indexer.NewMemoryIndexer()
	for n := 0; n < conjunctionNum; n++ {
		id := n + 1

		n1 := rand.Intn(5) + 1
		attrs := make([]*expr.Attribute, n1)
		for i := 0; i < n1; i++ {

			name := testAttrNames[rand.Intn(len(testAttrNames))]
			attrs[i] = &expr.Attribute{
				Name:     name,
				Values:   testAttrs[name],
				Contains: randBool(),
			}
		}

		k.Add(expr.NewConjunction(id, attrs))
	}
	k.Build()

	assignments := make([]expr.Assignment, assignmentNum)
	for i := 0; i < assignmentNum; i++ {
		assignments[i] = getTestAssignment(assignmentAvgSize, testAttrs, testAttrNames)
	}

	return k, assignments
}

var benchmarkResults []int

func Benchmark_Match_10000_20(b *testing.B) {

	k, assignments := getIndexerAndAssignment(10000, 10000, 20)

	b.ResetTimer()
	matcher := &matcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_100000_20(b *testing.B) {

	k, assignments := getIndexerAndAssignment(100000, 10000, 20)

	b.ResetTimer()
	matcher := &matcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_1000000_20(b *testing.B) {

	k, assignments := getIndexerAndAssignment(1000000, 10000, 20)

	b.ResetTimer()
	matcher := &matcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_10000_30(b *testing.B) {

	k, assignments := getIndexerAndAssignment(10000, 10000, 30)

	b.ResetTimer()
	matcher := &matcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_100000_30(b *testing.B) {

	k, assignments := getIndexerAndAssignment(100000, 10000, 30)

	b.ResetTimer()
	matcher := &matcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_1000000_30(b *testing.B) {

	k, assignments := getIndexerAndAssignment(1000000, 10000, 30)

	b.ResetTimer()
	matcher := &matcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_10000_40(b *testing.B) {

	k, assignments := getIndexerAndAssignment(10000, 10000, 40)

	b.ResetTimer()
	matcher := &matcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_100000_40(b *testing.B) {

	k, assignments := getIndexerAndAssignment(100000, 10000, 40)

	b.ResetTimer()
	matcher := &matcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_1000000_40(b *testing.B) {

	k, assignments := getIndexerAndAssignment(1000000, 10000, 40)

	b.ResetTimer()
	matcher := &matcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Test_kIndexTable_Match(t *testing.T) {

	k := indexer.NewMemoryIndexer()

	k.Add(expr.NewConjunction(
		1,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true},
			{Name: "state", Values: []string{"NY"}, Contains: true},
		},
	))

	k.Add(expr.NewConjunction(
		2,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true},
			{Name: "gender", Values: []string{"F"}, Contains: true},
		},
	))

	k.Add(expr.NewConjunction(
		3,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true},
			{Name: "gender", Values: []string{"M"}, Contains: true},
			{Name: "state", Values: []string{"CA"}, Contains: false},
		},
	))

	k.Add(expr.NewConjunction(
		4,
		[]*expr.Attribute{
			{Name: "state", Values: []string{"CA"}, Contains: true},
			{Name: "gender", Values: []string{"M"}, Contains: true},
		},
	))

	k.Add(expr.NewConjunction(
		5,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3", "4"}, Contains: true},
		},
	))

	k.Add(expr.NewConjunction(
		6,
		[]*expr.Attribute{
			{Name: "state", Values: []string{"CA", "NY"}, Contains: false},
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

	matcher := &matcher{}

	matched := matcher.Match(k, expr.Assignment{
		expr.Label{Name: "age", Value: "3"},
		expr.Label{Name: "state", Value: "CA"},
		expr.Label{Name: "gender", Value: "M"},
	})

	assert.ElementsMatch(t, []int{4, 5}, matched)
}
