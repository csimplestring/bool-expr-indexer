package main

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/matcher"
	"github.com/csimplestring/bool-expr-indexer/dnf/tools"
	"github.com/pkg/profile"
)

func main() {

	k, assignments := tools.GetPrefilledIndex(1000, 1000000, 10000, 10)

	tools.PrintMemUsage()

	defer profile.Start(profile.MemProfileHeap()).Stop()

	sm := matcher.New()
	for i := 0; i < 10000; i++ {
		sm.Match(k, assignments[i])
	}

	tools.PrintMemUsage()

	sm.Match(k, nil)
}
