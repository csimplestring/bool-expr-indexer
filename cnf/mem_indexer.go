package cnf

import "github.com/csimplestring/bool-expr-indexer/set"

// Key is the key representing an attribute, e.g., <age, 10>
type Key struct {
	Name  string
	Value string
	score int
}

func (k Key) hash() string {
	return k.Name + ":" + k.Value
}

type memoryIndexer struct {
	s map[string][]set.SortedSet
}

func (m *memoryIndexer) add(k string, entry *PostingEntry) {
	sortSets, exists := m.s[k]
	if !exists || sortSets == nil {
		sortSets = append(sortSets, set.NewSortedSet())
		m.s[k] = sortSets
	}

	inserted := false
	for _, set := range sortSets {
		if set.Add(entry) {
			inserted = true
			break
		}
	}

	if !inserted {
		m.s[k] = append(sortSets, set.NewSortedSet())
		m.add(k, entry)
	}
}
