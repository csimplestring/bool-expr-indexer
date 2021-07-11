package tools

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/api/dnf/indexer"
	"github.com/dchest/uniuri"
)

func randBool() bool {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(100)
	// it will return false in 1% possibility
	if n <= 1 {
		return false
	}
	return true
}

func PrintMemUsage() {
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

func getTestAttributes(num int) (map[string][]string, []string) {

	m := make(map[string][]string)

	names := make([]string, num)
	for i := 0; i < num; i++ {
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
		// we simulate a mis-match case here.
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

func GetPrefilledIndex(attributeNum, conjunctionNum, assignmentNum, assignmentAvgSize int) (indexer.Indexer, []expr.Assignment) {
	testAttrs, testAttrNames := getTestAttributes(attributeNum)

	k := indexer.NewMemIndexer()
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
