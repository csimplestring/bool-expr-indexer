package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer"
	"github.com/csimplestring/bool-expr-indexer/dnf/matcher/simple"
	"github.com/dchest/uniuri"
	"github.com/pkg/profile"
)

func randBool() bool {
	rand.Seed(time.Now().UnixNano())
	// n := rand.Intn(100)
	// if n <= 1 {
	// 	return false
	// }
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

		// var v string
		// n2 := rand.Intn(100)
		// if n2 >= len(values) {
		// 	v = ""
		// } else {
		// 	v = values[n2]
		// }
		labels[i] = expr.Label{
			Name:  name,
			Value: values[0],
		}
	}
	return labels
}

func main() {

	testAttrs, attrNames := getTestAttributes()

	k := indexer.NewMemoryIndexer(indexer.NewMapScorer())

	printMemUsage()

	for n := 0; n < 1000000; n++ {
		id := n + 1

		n1 := rand.Intn(2) + 1
		attrs := make([]*expr.Attribute, n1)
		for i := 0; i < n1; i++ {

			name := attrNames[rand.Intn(len(attrNames))]
			attrs[i] = &expr.Attribute{
				Name:     name,
				Values:   testAttrs[name],
				Contains: randBool(),
			}
		}

		conj := expr.NewConjunction(id, attrs)

		k.Add(conj)
		// b, _ := json.Marshal(conj)
		// fmt.Println(string(b))
	}

	k.Build()

	runtime.GC()

	fmt.Println()

	assignments := make([]expr.Assignment, 10000)
	for i := 0; i < 10000; i++ {
		assignments[i] = getTestAssignment(5, testAttrs, attrNames)
	}

	defer profile.Start(profile.CPUProfile).Stop()

	sm := simple.New()
	for i := 0; i < 10000; i++ {
		sm.Match(k, assignments[i])
	}
}
