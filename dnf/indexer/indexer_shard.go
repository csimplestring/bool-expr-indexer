package indexer

import (
	"os"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
)

const zeroKey string = ":"

type shard interface {
	Load(f os.File)
	Get(name string, value string) *Record
}

// indexShard stores the posting indexes for all the conjunctions with the same size.
type indexShard struct {
	conjunctionSize int
	zeroKey         string
	invertedMap     map[string]*Record
}

// newIndexShard creates a new indexShard.
func newIndexShard(ksize int) *indexShard {

	return &indexShard{
		zeroKey:         zeroKey,
		conjunctionSize: ksize,
		invertedMap:     make(map[string]*Record),
	}
}

func (m *indexShard) hashKey(name string, value string) string {
	return name + ":" + value
}

func (m *indexShard) Build() error {

	for _, r := range m.invertedMap {
		r.PostingList.Sort()
		r.compact()
	}
	return nil
}

func (m *indexShard) Get(name string, value string) *Record {
	v, ok := m.invertedMap[m.hashKey(name, value)]
	if !ok {
		return nil
	}

	return v
}

func (m *indexShard) createIfAbsent(hash string, name, value string) *Record {

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

			// todo: make it concurrent saf
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
