package indexer

import (
	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/csimplestring/bool-expr-indexer/dnf/indexer/posting"
)

type indexShard struct {
	conjunctionSize int
	zeroKey         string
	invertedMap     map[string]*Record
	//hash            maphash.Hash
}

func newIndexShard(ksize int) *indexShard {
	// init a hasher
	// var hasher maphash.Hash
	// hasher.Reset()
	// hasher.WriteString("")
	// hasher.WriteString("")
	// zeroKey := hasher.Sum64()

	return &indexShard{
		zeroKey:         ":",
		conjunctionSize: ksize,
		//hash:            hasher,
		invertedMap: make(map[string]*Record),
	}
}

func (m *indexShard) hashKey(name string, value string) string {
	// m.hash.Reset()
	// m.hash.WriteString(name)
	// m.hash.WriteString(value)
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
	h := m.hashKey(name, value)
	return m.invertedMap[h]
}

func (m *indexShard) createIfAbsent(hash string, name, value string) *Record {
	v := m.invertedMap[hash]

	if v == nil {
		r := &Record{
			PostingList: make(posting.List, 0, 64),
			Key:         name,
			Value:       value,
		}
		m.invertedMap[hash] = r
		return r
	}
	return v
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
