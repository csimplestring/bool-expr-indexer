package dnf

import (
	"github.com/csimplestring/bool-expr-indexer/set"
)

// Matcher finds the matched conjunction ids
type Matcher interface {
	Match(*kIndexTable, Assignment) []int64
}

type matcher struct {
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *matcher) Match(k *kIndexTable, labels Assignment) []int64 {
	results := set.Int64HashSet()

	n := min(len(labels), k.maxKSize)

	for i := n; i >= 0; i-- {
		pLists := newPostingLists(k.GetPostingLists(i, labels))

		K := i
		if K == 0 {
			K = 1
		}
		if pLists.len() < K {
			continue
		}

		pLists.sortByCurrent()
		for pLists.c[K-1].current() != eolItem {
			var nextID int64

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
