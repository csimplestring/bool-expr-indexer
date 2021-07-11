package matcher

import (
	"sort"

	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
)

type ranker struct {
	assignment    expr.Assignment
	sWeightLookup map[string]int
}

func newRanker(assignment expr.Assignment) *ranker {
	u := &ranker{
		assignment: assignment,
	}

	sWeightLookup := make(map[string]int)
	for _, label := range assignment {
		sWeightLookup[u.formatKey(label.Name, label.Value)] = label.Weight
	}
	u.sWeightLookup = sWeightLookup

	return u
}

func (u *ranker) formatKey(name, value string) string {
	return name + ":" + value
}

// lists must be sorted
func (u *ranker) calculateEntryUB(K int, lists postingLists) int {

	score := 0
	for i := 0; i <= K && i < lists.Len(); i++ {
		name := lists[i].name
		value := lists[i].value
		if w, exists := u.sWeightLookup[u.formatKey(name, value)]; exists {
			score += w * lists[i].ub
		}
	}
	return score
}

func (u *ranker) calculateListUB(topN int, lists postingLists) int {

	scores := make([]int, lists.Len())
	for i, list := range lists {
		if w, exists := u.sWeightLookup[u.formatKey(list.name, list.value)]; exists {
			scores[i] = w * list.ub
		} else {
			scores[i] = 0
		}
	}
	// get topN
	sort.Ints(scores)
	totalUB := 0
	for i := 0; i < topN && i < len(scores); i++ {
		totalUB += scores[i]
	}

	return totalUB
}
