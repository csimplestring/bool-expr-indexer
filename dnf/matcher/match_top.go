package matcher

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
	"github.com/csimplestring/bool-expr-indexer/dnf/matcher/pq"
	"github.com/csimplestring/bool-expr-indexer/dnf/scorer"
)

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

		listUB := ubCalculator.calculateListUB(K, pLists)
		if q.Len() == topN && q.PeekMin().Priority() > listUB {
			continue
		}

		pLists.sortByCurrent()
		for pLists[K-1].current() != posting.EOL {
			conjunctionUB := ubCalculator.calculateEntryUB(K-1, pLists)

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

func extractIDs(q pq.MinMaxPriorityQueue) []int {
	ids := make([]int, q.Len())
	i := 0
	for q.Len() != 0 {
		ids[i] = q.PopMin().(*pq.IntItem).Value().(int)
		i++
	}
	return ids
}
