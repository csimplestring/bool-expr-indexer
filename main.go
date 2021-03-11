package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"runtime"
// 	"time"

// 	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
// 	"github.com/csimplestring/bool-expr-indexer/dnf/indexer"
// 	"github.com/dchest/uniuri"
// 	"github.com/pkg/profile"
// )

// func randBool() bool {
// 	rand.Seed(time.Now().UnixNano())
// 	return rand.Intn(2) == 1
// }

// func printMemUsage() {
// 	bToMb := func(b uint64) uint64 {
// 		return b / 1024 / 1024
// 	}
// 	var m runtime.MemStats
// 	runtime.ReadMemStats(&m)
// 	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
// 	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
// 	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
// 	fmt.Printf("\tNumGC = %v\n", m.NumGC)
// }

// func getTestAttributes() map[string][]string {

// 	m := make(map[string][]string)

// 	for i := 0; i < 1000; i++ {
// 		nameLen := rand.Intn(8) + 1
// 		name := uniuri.NewLen(nameLen)

// 		valuesLen := rand.Intn(20) + 1
// 		values := make([]string, valuesLen)
// 		for j := 0; j < valuesLen; j++ {
// 			values[j] = uniuri.NewLen(valuesLen)
// 		}

// 		m[name] = values
// 	}

// 	return m
// }

// func main() {
// 	defer profile.Start(profile.MemProfile).Stop()

// 	metastorer := expr.NewAttributeMetadataStorer()

// 	testAttrs := getTestAttributes()
// 	i := uint32(1)
// 	var attrNames []string
// 	for name, values := range testAttrs {
// 		metastorer.AddNameIDMapping(name, i)
// 		i++
// 		attrNames = append(attrNames, name)

// 		for j, v := range values {
// 			metastorer.AddValueIDMapping(name, v, uint32(j))
// 		}
// 	}

// 	k := indexer.NewMemoryIndexer(metastorer)

// 	printMemUsage()

// 	for n := 0; n < 1000000; n++ {
// 		id := n + 1

// 		var attrs []expr.Attribute
// 		for i := 0; i < rand.Intn(5)+1; i++ {
// 			name := attrNames[rand.Intn(len(attrNames))]
// 			attr, err := metastorer.NewAttribute(name, testAttrs[name], randBool())

// 			if err != nil {
// 				panic(err)
// 			}
// 			attrs = append(attrs, attr)
// 		}

// 		k.Add(expr.NewConjunction(id, attrs))
// 	}

// 	k.Build()

// 	printMemUsage()
// }
