package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/api/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/api/dnf/indexer/posting"
)

const zeroKey string = ":"

// shard is a sub-map of indexer, which only stores the conjunctions with same size.
type shard interface {
	Get(name string, value string) *Record
	Build() error
	Add(c *expr.Conjunction) error
}

// mapShard implements shard, stores the posting indexes for all the conjunctions with the same size.
type mapShard struct {
	conjunctionSize int
	zeroKey         string
	invertedMap     map[string]*Record
}

// newMapShard creates a new mapShard.
func newMapShard(ksize int) *mapShard {

	return &mapShard{
		zeroKey:         zeroKey,
		conjunctionSize: ksize,
		invertedMap:     make(map[string]*Record),
	}
}

// hashKey concats the name:value as hash key.
func (m *mapShard) hashKey(name string, value string) string {
	return name + ":" + value
}

// Build finalise the posting lists in shard, sort and compact each post list.
func (m *mapShard) Build() error {

	for _, r := range m.invertedMap {
		r.PostingList.Sort()
		r.compact()
	}
	return nil
}

// Get returns Record based on name:value
func (m *mapShard) Get(name string, value string) *Record {
	v, ok := m.invertedMap[m.hashKey(name, value)]
	if !ok {
		return nil
	}

	return v
}

// createIfAbsent
func (m *mapShard) createIfAbsent(hash string, name, value string) *Record {

	if v, found := m.invertedMap[hash]; found {
		return v
	}

	m.invertedMap[hash] = &Record{
		PostingList: make(posting.List, 0, 64),
		Key:         name,
		Value:       value,
	}

	return m.invertedMap[hash]
}

// Add conjunction into a shard.
func (m *mapShard) Add(c *expr.Conjunction) error {

	for _, attr := range c.Attributes {
		for i, value := range attr.Values {

			hash := m.hashKey(attr.Name, value)

			r := m.createIfAbsent(hash, attr.Name, value)

			score := uint32(0)
			if len(attr.Weights) != 0 {
				score = attr.Weights[i]
			}

			entry, err := posting.NewEntryInt32(uint32(c.ID), attr.Contains, score)
			if err != nil {
				return err
			}

			r.append(entry)

			m.invertedMap[hash] = r
		}
	}

	if c.GetKSize() == 0 {
		r := m.createIfAbsent(m.zeroKey, "", "")

		entry, err := posting.NewEntryInt32(uint32(c.ID), true, 0)
		if err != nil {
			return err
		}
		r.append(entry)
		m.invertedMap[m.zeroKey] = r
	}

	return nil
}
