package matcher

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
)

type AllMatcher interface {
	Match(indexer indexer.Indexer, assignment expr.Assignment) []int
}

func New() AllMatcher {
	return &allMatcher{}
}

type allMatcher struct{}

func (m *allMatcher) Match(indexer indexer.Indexer, assignment expr.Assignment) []int {
	results := make([]int, 0, 1024)

	n := min(len(assignment), indexer.MaxKSize())

	for i := n; i >= 0; i-- {
		pLists := newPostingLists(indexer.Get(i, assignment), nil)

		K := i
		if K == 0 {
			K = 1
		}
		if pLists.Len() < K {
			continue
		}

		pLists.sortByCurrent()
		for pLists[K-1].current() != posting.EOL {
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

	return results
}
