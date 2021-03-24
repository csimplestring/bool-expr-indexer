package matcher

import (
	"sort"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
	"github.com/csimplestring/bool-expr-indexer/dnf/matcher/pq"
	"github.com/csimplestring/bool-expr-indexer/dnf/scorer"
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
func (u *ranker) conjunctionUB(K int, lists postingLists) int {

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

func (u *ranker) listUB(topN int, lists postingLists) int {

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

func extractIDs(q pq.MinMaxPriorityQueue) []int {
	ids := make([]int, q.Len())
	i := 0
	for q.Len() != 0 {
		ids[i] = q.PopMin().(*pq.IntItem).Value().(int)
		i++
	}
	return ids
}

type TopNMatcher interface {
	MatchTopN(topN int, indexer indexer.Indexer, assignment expr.Assignment) []int
}

func NewTopN(s scorer.Scorer) TopNMatcher {
	return &topNMatcher{
		scorer: s,
	}
}

type topNMatcher struct {
	scorer scorer.Scorer
}

func (m *topNMatcher) MatchTopN(topN int, indexer indexer.Indexer, assignment expr.Assignment) []int {
	results := make([]int, 0, 1024)
	ubCalculator := newRanker(assignment)
	q := pq.New(topN)

	n := min(len(assignment), indexer.MaxKSize())

	for i := n; i >= 0; i-- {
		pLists := newPostingLists(indexer.Get(i, assignment), m.scorer)

		K := i
		if K == 0 {
			K = 1
		}
		if pLists.Len() < K {
			continue
		}

		listUB := ubCalculator.listUB(K, pLists)
		if q.Len() == topN && q.PeekMin().Priority() > listUB {
			continue
		}

		pLists.sortByCurrent()
		for pLists[K-1].current() != posting.EOL {
			conjunctionUB := ubCalculator.conjunctionUB(K-1, pLists)

			if q.Len() == topN && q.PeekMin().Priority() > conjunctionUB {
				nextID := pLists[K-1].current().CID() + 1
				for L := 0; L <= K-1; L++ {
					pLists[L].skipTo(nextID)
				}
				pLists.sortByCurrent()
				continue
			}

			var nextID uint32

			if pLists[0].current().CID() == pLists[K-1].current().CID() {

				if pLists[0].current().Contains() == false {
					rejectID := pLists[0].current().CID()
					for L := K; L <= pLists.Len()-1; L++ {
						if pLists[L].current().CID() == rejectID {
							pLists[L].skipTo(rejectID + 1)
						} else {
							break
						}
					}

				} else {
					q.Push(&pq.IntItem{Val: int(pLists[K-1].current().CID()), Prior: conjunctionUB})
					results = append(results, int(pLists[K-1].current().CID()))
				}

				nextID = pLists[K-1].current().CID() + 1
			} else {
				nextID = pLists[K-1].current().CID()

			}

			for L := 0; L <= K-1; L++ {
				pLists[L].skipTo(nextID)
			}
			pLists.sortByCurrent()
		}

	}

	return extractIDs(q)
}
