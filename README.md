# Boolean Expression Indexer Go library

A Go implementation of the core algorithm in paper <[Indexing Boolean Expression](https://theory.stanford.edu/~sergei/papers/vldb09-indexing.pdf)>, which already supports the following features mentioned in paper:

- DNF algorithm
- Simple match
- TopN match (ranking based)

# Use Scenario

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

Total number of expressions, value size                              

Benchmark_Match_10000_20-12                       106482             10142 ns/op

Benchmark_Match_100000_20-12                       70419             14749 ns/op

Benchmark_Match_1000000_20-12                      14438             87884 ns/op

Benchmark_Match_10000_30-12                        80452             15832 ns/op

Benchmark_Match_100000_30-12                       61867             20770 ns/op

Benchmark_Match_1000000_30-12                      10000            103594 ns/op

Benchmark_Match_10000_40-12                        61221             20778 ns/op

Benchmark_Match_100000_40-12                       49161             25129 ns/op

Benchmark_Match_1000000_40-12                      10000            110132 ns/op

Memory usage: 1 million of expressions are indexed and it takes 100 MB on average.

## usage

see matcher_test.go 

 

## Roadmap

This library only implements the core DNF algorithm in that paper. However, it is more useful to use it to build up a production-ready Ad Server. Currently, the indexer does not support online update(all the expressions shall be indexed once in the beginning). Also, the speed will be slow if the expressions number increases. Therefore I plan to support more advanced features soon

- online index update
- metrics monitoring
- index partitioning
- canary rollout deployment
- HA support
- http/gRPC transport
- web UI

PR issues are welcome