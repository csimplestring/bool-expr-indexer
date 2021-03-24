package matcher

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer"
	"github.com/csimplestring/bool-expr-indexer/dnf/scorer"
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
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_100000_20(b *testing.B) {

	k, assignments := getIndexerAndAssignment(100000, 10000, 20)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_1000000_20(b *testing.B) {

	k, assignments := getIndexerAndAssignment(1000000, 10000, 20)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_10000_30(b *testing.B) {

	k, assignments := getIndexerAndAssignment(10000, 10000, 30)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_100000_30(b *testing.B) {

	k, assignments := getIndexerAndAssignment(100000, 10000, 30)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_1000000_30(b *testing.B) {

	k, assignments := getIndexerAndAssignment(1000000, 10000, 30)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_10000_40(b *testing.B) {

	k, assignments := getIndexerAndAssignment(10000, 10000, 40)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_100000_40(b *testing.B) {

	k, assignments := getIndexerAndAssignment(100000, 10000, 40)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_1000000_40(b *testing.B) {

	k, assignments := getIndexerAndAssignment(1000000, 10000, 40)

	b.ResetTimer()
	matcher := &allMatcher{}
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
			{Name: "age", Values: []string{"3"}, Contains: true, Weights: []uint32{1}},
			{Name: "state", Values: []string{"NY"}, Contains: true, Weights: []uint32{40}},
		},
	))

	k.Add(expr.NewConjunction(
		2,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true, Weights: []uint32{1}},
			{Name: "gender", Values: []string{"F"}, Contains: true, Weights: []uint32{3}},
		},
	))

	k.Add(expr.NewConjunction(
		3,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true, Weights: []uint32{2}},
			{Name: "gender", Values: []string{"M"}, Contains: true, Weights: []uint32{5}},
			{Name: "state", Values: []string{"CA"}, Contains: false, Weights: []uint32{0}},
		},
	))

	k.Add(expr.NewConjunction(
		4,
		[]*expr.Attribute{
			{Name: "state", Values: []string{"CA"}, Contains: true, Weights: []uint32{15}},
			{Name: "gender", Values: []string{"M"}, Contains: true, Weights: []uint32{9}},
		},
	))

	k.Add(expr.NewConjunction(
		5,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3", "4"}, Contains: true, Weights: []uint32{1, 5}},
		},
	))

	k.Add(expr.NewConjunction(
		6,
		[]*expr.Attribute{
			{Name: "state", Values: []string{"CA", "NY"}, Contains: false, Weights: []uint32{0, 0}},
		},
	))

	k.Build()

	scoreMap := scorer.NewMapScorer()
	scoreMap.SetUB("state", "CA", 2)
	scoreMap.SetUB("state", "NY", 5)
	scoreMap.SetUB("age", "3", 1)
	scoreMap.SetUB("age", "4", 3)
	scoreMap.SetUB("gender", "F", 2)
	scoreMap.SetUB("gender", "M", 1)

	assert.Equal(t, 2, k.MaxKSize())
	matcher := &allMatcher{}

	matched := matcher.Match(k, expr.Assignment{
		expr.Label{Name: "age", Value: "3"},
		expr.Label{Name: "state", Value: "CA"},
		expr.Label{Name: "gender", Value: "M"},
	})

	assert.ElementsMatch(t, []int{4, 5}, matched)

	matched = matcher.Match(k, expr.Assignment{
		expr.Label{Name: "age", Value: "3"},
		expr.Label{Name: "state", Value: "NY"},
		expr.Label{Name: "gender", Value: "F"},
	})

	assert.ElementsMatch(t, []int{1, 2, 5}, matched)

	topNMatcher := NewTopN(scoreMap)
	matched = topNMatcher.MatchTopN(1, k, expr.Assignment{
		expr.Label{Name: "age", Value: "3", Weight: 8},
		expr.Label{Name: "state", Value: "NY", Weight: 10},
		expr.Label{Name: "gender", Value: "F", Weight: 9},
	})
	assert.ElementsMatch(t, []int{1}, matched)
}
