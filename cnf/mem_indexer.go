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
	sortedSetMap map[string][]set.SortedSet
	pListsMap    map[string]PostingLists
}

func (m *memoryIndexer) buildPostingList(s set.SortedSet) PostingList {
	sl := s.ToSlice()
	p := make(PostingList, len(sl))
	for i, item := range sl {
		p[i] = item.(*PostingEntry)
	}

	return p
}

func (m *memoryIndexer) Build() {
	m.pListsMap = make(map[string]PostingLists, len(m.sortedSetMap))
	for key, sets := range m.sortedSetMap {
		p := make(PostingLists, len(sets))
		for i, set := range sets {
			p[i] = m.buildPostingList(set)
		}
		m.pListsMap[key] = p
	}
}

func (m *memoryIndexer) add(k string, entry *PostingEntry) {
	sortSets, exists := m.sortedSetMap[k]
	if !exists || sortSets == nil {
		sortSets = append(sortSets, set.NewSortedSet())
		m.sortedSetMap[k] = sortSets
	}

	inserted := false
	for _, set := range sortSets {
		if set.Add(entry) {
			inserted = true
			break
		}
	}

	if !inserted {
		m.sortedSetMap[k] = append(sortSets, set.NewSortedSet())
		m.add(k, entry)
	}
}
