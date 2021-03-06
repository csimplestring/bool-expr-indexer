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
	Get(k *Key) *PostingList
	Build() error
}

// Key is the key representing an attribute, e.g., <age, 10>
type Key struct {
	Name  uint32
	Value uint32
	score int
}
