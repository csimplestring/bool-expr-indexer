package hash

import (
	"github.com/cespare/xxhash/v2"
)

type Hasher interface {
	WriteString(string) (int, error)
	Sum64() uint64
	Reset()
}

func NewXxHash() Hasher {
	return xxhash.New()
}

func ZeroKey(h Hasher) uint64 {
	h.Reset()
	h.WriteString("")
	h.WriteString("")
	return h.Sum64()
}
