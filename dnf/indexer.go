package dnf

// KSizeIndexer shards the Indexer by conjunction size.
type KSizeIndexer interface {
	Build() error
	MaxKSize() int
	Get(size int) Indexer
}

// Indexer actually stores the reverted index: key -> posting list
type Indexer interface {
	Add(c *Conjunction) error
	Get(k *key) *PostingList
	Build() error
}

// key is the key representing an attribute, e.g., <age, 10>
type key struct {
	Name  uint32
	Value uint32
	score int
}

// zKey is zero key placeholder
var zKey *key = &key{
	Name:  0,
	Value: 0,
	score: 0,
}
