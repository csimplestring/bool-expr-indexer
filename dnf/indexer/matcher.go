package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/set"
)

// Matcher finds the matched conjunction ids
type Matcher interface {
	Match(expr.Assignment) []int
}

// NewMatcher creates a new matcher
func NewMatcher(k *kIndexTable) Matcher {
	return &matcher{
		kIndexTable: k,
	}
}

// matcher implements the Matcher interface
type matcher struct {
	kIndexTable *kIndexTable
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// Match finds the matched conjunctions given an assignment.
func (m *matcher) Match(assignment expr.Assignment) []int {
	results := set.IntHashSet()
	k := m.kIndexTable

	n := min(len(assignment), k.maxKSize)

	for i := n; i >= 0; i-- {
		pLists := newPostingLists(k.GetPostingLists(i, assignment))

		K := i
		if K == 0 {
			K = 1
		}
		if pLists.len() < K {
			continue
		}

		pLists.sortByCurrent()
		for pLists.c[K-1].current() != eolItem {
			var nextID int

			if pLists.c[0].current().CID == pLists.c[K-1].current().CID {

				if pLists.c[0].current().Contains == false {
					rejectID := pLists.c[0].current().CID
					for L := K; L <= pLists.len()-1; L++ {
						if pLists.c[L].current().CID == rejectID {
							pLists.c[L].skipTo(rejectID + 1)
						} else {
							break
						}
					}

				} else {
					results.Add(pLists.c[K-1].current().CID)
				}

				nextID = pLists.c[K-1].current().CID + 1
			} else {
				nextID = pLists.c[K-1].current().CID

			}

			for L := 0; L <= K-1; L++ {
				pLists.c[L].skipTo(nextID)
			}
			pLists.sortByCurrent()
		}

	}

	return results.ToSlice()
}
