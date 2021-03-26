package indexer

import (
	"github.com/cornelk/hashmap"
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
)

type indexShard struct {
	conjunctionSize int
	zeroKey         string
	invertedMap     *hashmap.HashMap
}

func newIndexShard(ksize int) *indexShard {

	return &indexShard{
		zeroKey:         ":",
		conjunctionSize: ksize,
		invertedMap:     &hashmap.HashMap{},
	}
}

func (m *indexShard) hashKey(name string, value string) string {
	return name + ":" + value
}

func (m *indexShard) Build() error {

	for next := range m.invertedMap.Iter() {
		r := next.Value.(*Record)
		r.PostingList.Sort()
		r.compact()
	}
	return nil
}

func (m *indexShard) Get(name string, value string) *Record {
	v, ok := m.invertedMap.GetStringKey(m.hashKey(name, value))
	if !ok {
		return nil
	}

	return v.(*Record)
}

func (m *indexShard) createIfAbsent(hash string, name, value string) *Record {
	key := m.hashKey(name, value)
	v, _ := m.invertedMap.GetOrInsert(key, &Record{
		PostingList: make(posting.List, 0, 64),
		Key:         name,
		Value:       value,
	})

	return v.(*Record)
}

func (m *indexShard) Add(c *expr.Conjunction) error {

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

			// todo: make it concurrent safe
			m.invertedMap.Set(hash, r)
		}
	}

	if c.GetKSize() == 0 {
		r := m.createIfAbsent(m.zeroKey, "", "")

		entry, err := posting.NewEntryInt32(uint32(c.ID), true, 0)
		if err != nil {
			return err
		}
		r.append(entry)
		m.invertedMap.Set(m.zeroKey, r)
	}

	return nil
}
