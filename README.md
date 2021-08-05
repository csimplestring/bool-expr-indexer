# Boolean Expression Indexer Go library

![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/csimplestring/bool-expr-indexer/Go/master?style=for-the-badge)
[![GitHub issues](https://img.shields.io/github/issues/csimplestring/bool-expr-indexer?style=for-the-badge)](https://github.com/csimplestring/bool-expr-indexer/issues)
[![GitHub stars](https://img.shields.io/github/stars/csimplestring/bool-expr-indexer?style=for-the-badge)](https://github.com/csimplestring/bool-expr-indexer/stargazers)
[![GitHub license](https://img.shields.io/github/license/csimplestring/bool-expr-indexer?style=for-the-badge)](https://github.com/csimplestring/bool-expr-indexer/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/csimplestring/bool-expr-indexer?style=for-the-badge)](https://goreportcard.com/report/github.com/csimplestring/bool-expr-indexer)



A Go implementation of the core algorithm in paper <[Indexing Boolean Expression](https://theory.stanford.edu/~sergei/papers/vldb09-indexing.pdf)>, which already supports the following features mentioned in paper:

- DNF algorithm
- Simple match
- TopN match (ranking based)
- Online update (Copy-On-Write)


## usage 

``` Go

import 
	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/api/dnf/indexer"
	"github.com/stretchr/testify/assert"
)

    var conjunctions []*expr.Conjunction

	conjunctions = append(conjunctions, expr.NewConjunction(
		1,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true, Weights: []uint32{1}},
			{Name: "state", Values: []string{"NY"}, Contains: true, Weights: []uint32{40}},
		},
	))

	conjunctions = append(conjunctions, expr.NewConjunction(
		2,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true, Weights: []uint32{1}},
			{Name: "gender", Values: []string{"F"}, Contains: true, Weights: []uint32{3}},
		},
	))

	conjunctions = append(conjunctions, expr.NewConjunction(
		3,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true, Weights: []uint32{2}},
			{Name: "gender", Values: []string{"M"}, Contains: true, Weights: []uint32{5}},
			{Name: "state", Values: []string{"CA"}, Contains: false, Weights: []uint32{0}},
		},
	))

	conjunctions = append(conjunctions, expr.NewConjunction(
		4,
		[]*expr.Attribute{
			{Name: "state", Values: []string{"CA"}, Contains: true, Weights: []uint32{15}},
			{Name: "gender", Values: []string{"M"}, Contains: true, Weights: []uint32{9}},
		},
	))

	conjunctions = append(conjunctions, expr.NewConjunction(
		5,
		[]*expr.Attribute{
			{Name: "age", Values: []string{"3", "4"}, Contains: true, Weights: []uint32{1, 5}},
		},
	))

	conjunctions = append(conjunctions, expr.NewConjunction(
		6,
		[]*expr.Attribute{
			{Name: "state", Values: []string{"CA", "NY"}, Contains: false, Weights: []uint32{0, 0}},
		},
	))

	// read only indexer, best performance
	k, err := indexer.NewMemReadOnlyIndexer(conjunctions)
	assert.NoError(t, err)
	k.Build()

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

	// copy on write indexer, you can ADD or DELETE new conjunctions
	// the refresh interval in second, to control how often to swap the old and new index
	// usually set it up to 15 minutes, which is ok in most cases.
	freshInterval := 5
	k, err := indexer.NewCopyOnWriteIndexer(conjunctions, freshInterval)
	assert.NoError(t, err)

	k.Delete(4)
	
	matched = matcher.Match(k, expr.Assignment{
		expr.Label{Name: "age", Value: "3"},
		expr.Label{Name: "state", Value: "CA"},
		expr.Label{Name: "gender", Value: "M"},
	})
	assert.ElementsMatch(t, []int{4, 5}, matched)

	// after 5 seconds, the modifications will be applied.
	time.Sleep(5 * time.Second)
	matched = matcher.Match(k, expr.Assignment{
		expr.Label{Name: "age", Value: "3"},
		expr.Label{Name: "state", Value: "CA"},
		expr.Label{Name: "gender", Value: "M"},
	})
	assert.ElementsMatch(t, []int{5}, matched)
```

see matcher_test.go 

## Use Scenario

It is an important component in online computational advertising system. As mentioned in the paper, the original motivation of this library is to provide an efficient way to select ads based on pre-defined targeting condition. Given a user with some labels, it finds the best-matched ads. 

In a typical advertising management system, the following hierarchy models are used:

    Advertisers
    |__ LineItem 
        |__ InsertionOrder 
            |__ Campaign
                |__ Creative Banners With Targeting Condition

### Ad selector in RTB

In the RTB or similar environment, when a bidding request comes, usually it comes with some user-profile context(in the paper, it is called ***assignment***). The Ad server then needs to quickly find the matched creative banners that match the targeting condition(in the paper, it is expressed as ***conjunctions***).  A naive way is to iterate all the banners and compare the user-profile with targeting condition one-by-one, which is very slow when there are millions of banners in the system, however the whole RTB workflow must be done with 100 ms.

### User segmentation in DMP

In a DMP environment, all the collected user data shall be processed, enriched and aggregated in nearly real-time. One key problem is to quickly identify a user belongs to which group, based on pre-defined group condition. 

This library is developed for the above scenarios. The following features are supported

- DNF algorithm

    For example: a targeting condition can be expressed as: (age = 30 AND city = A) OR (gender = male AND age = 20 AND region = B) etc

- Simple match

    It returns all the matched expression ids

- TopN match (ranking based)

    It returns the TopN matched expression ids, based on an adopted WAND algorithm. Note that the caller has to provide the upper bound and score of each conjunction, usually can be calculated in a Spark job. 

## Benchmark

Total number of expressions - size of key | operations within 1s |  op per ns |                         
|---:|---:|---:|
| Benchmark_Match_10000_20-12         |              106482  |            10142 ns/op |
| Benchmark_Match_100000_20-12          |             70419     |        14749 ns/op|
|Benchmark_Match_1000000_20-12          |            14438  |           87884 ns/op |
| Benchmark_Match_10000_30-12       |                 80452         |    15832 ns/op|
| Benchmark_Match_100000_30-12       |                61867       |      20770 ns/op |
| Benchmark_Match_1000000_30-12         |             10000      |      103594 ns/op |
|Benchmark_Match_10000_40-12     |                   61221         |    20778 ns/op |
| Benchmark_Match_100000_40-12  |                     49161        |     25129 ns/op |
| Benchmark_Match_1000000_40-12       |               10000        |    110132 ns/op |

Memory usage: 1 million of expressions are indexed and it takes 100 MB on average.


## Roadmap

This library only implements the core DNF algorithm in that paper. However, it is more useful to use it to build up a production-ready Ad Server. Currently, the indexer does not support online update(all the expressions shall be indexed once in the beginning). Also, the speed will be slow if the expressions number increases. Therefore I plan to support more advanced features soon

- metrics monitoring
- index partitioning
- canary rollout deployment
- HA support
- http/gRPC transport
- web UI

PR issues are welcome
