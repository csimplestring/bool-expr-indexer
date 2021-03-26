package matcher

import (
	"math/rand"
	"testing"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer"
	"github.com/csimplestring/bool-expr-indexer/dnf/scorer"
	"github.com/csimplestring/bool-expr-indexer/dnf/tools"
	"github.com/stretchr/testify/assert"
)

var benchmarkResults []int

func Benchmark_Match_10000_20(b *testing.B) {

	k, assignments := tools.GetPrefilledIndex(1000, 10000, 10000, 20)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_100000_20(b *testing.B) {

	k, assignments := tools.GetPrefilledIndex(1000, 100000, 10000, 20)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_1000000_20(b *testing.B) {

	k, assignments := tools.GetPrefilledIndex(1000, 1000000, 10000, 20)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_10000_30(b *testing.B) {

	k, assignments := tools.GetPrefilledIndex(1000, 10000, 10000, 30)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_100000_30(b *testing.B) {

	k, assignments := tools.GetPrefilledIndex(1000, 100000, 10000, 30)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_1000000_30(b *testing.B) {

	k, assignments := tools.GetPrefilledIndex(1000, 1000000, 10000, 30)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_10000_40(b *testing.B) {

	k, assignments := tools.GetPrefilledIndex(1000, 10000, 10000, 40)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_100000_40(b *testing.B) {

	k, assignments := tools.GetPrefilledIndex(1000, 100000, 10000, 40)

	b.ResetTimer()
	matcher := &allMatcher{}
	var result []int
	for i := 0; i < b.N; i++ {
		result = matcher.Match(k, assignments[rand.Intn(10000)])
	}
	benchmarkResults = result
}

func Benchmark_Match_1000000_40(b *testing.B) {

	k, assignments := tools.GetPrefilledIndex(1000, 1000000, 10000, 40)

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