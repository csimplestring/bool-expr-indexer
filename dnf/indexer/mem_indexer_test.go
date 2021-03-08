package indexer

import (
	"reflect"
	"testing"

	"github.com/csimplestring/bool-expr-indexer/dnf/expr"
	"github.com/stretchr/testify/assert"
)

func Test_indexShard_toKeys(t *testing.T) {
	m := &indexShard{}

	k0 := m.toKeys(nil)
	assert.Nil(t, k0)

	k1 := m.toKeys(&expr.Attribute{Name: 1, Values: []uint32{1}})
	assert.Equal(t, []*key{{Name: 1, Value: 1}}, k1)

	k2 := m.toKeys(&expr.Attribute{Name: 1, Values: []uint32{1, 2}})
	assert.ElementsMatch(t, []*key{{Name: 1, Value: 1}, {Name: 1, Value: 2}}, k2)
}

func Test_indexShard_hashKey(t *testing.T) {

}

func Test_indexShard_build(t *testing.T) {
	type fields struct {
		invertedMap   map[uint64]*PostingList
		attributeMeta expr.AttributeMetadataStorer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &indexShard{
				invertedMap:   tt.fields.invertedMap,
				attributeMeta: tt.fields.attributeMeta,
			}
			if err := m.Build(); (err != nil) != tt.wantErr {
				t.Errorf("indexShard.build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_indexShard_getPostingLists(t *testing.T) {

}

func Test_indexShard_createIfAbsent(t *testing.T) {
	type fields struct {
		invertedMap   map[uint64]*PostingList
		attributeMeta expr.AttributeMetadataStorer
	}
	type args struct {
		hash uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *PostingList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &indexShard{
				invertedMap:   tt.fields.invertedMap,
				attributeMeta: tt.fields.attributeMeta,
			}
			if got := m.createIfAbsent(tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("indexShard.createIfAbsent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexShard_get(t *testing.T) {
	type fields struct {
		invertedMap   map[uint64]*PostingList
		attributeMeta expr.AttributeMetadataStorer
	}
	type args struct {
		hash uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *PostingList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &indexShard{
				invertedMap:   tt.fields.invertedMap,
				attributeMeta: tt.fields.attributeMeta,
			}
			if got := m.get(tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("indexShard.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexShard_put(t *testing.T) {
	type fields struct {
		invertedMap   map[uint64]*PostingList
		attributeMeta expr.AttributeMetadataStorer
	}
	type args struct {
		hash uint64
		p    *PostingList
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &indexShard{
				invertedMap:   tt.fields.invertedMap,
				attributeMeta: tt.fields.attributeMeta,
			}
			m.put(tt.args.hash, tt.args.p)
		})
	}
}

func Test_indexShard_Add(t *testing.T) {

}

func TestNewMemoryIndexer(t *testing.T) {
	type args struct {
		attributeMeta expr.AttributeMetadataStorer
	}
	tests := []struct {
		name string
		args args
		want Indexer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryIndexer(tt.args.attributeMeta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoryIndexer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryIndex_Add(t *testing.T) {

}

func Test_memoryIndex_Build(t *testing.T) {
	type fields struct {
		maxKSize      int
		attributeMeta expr.AttributeMetadataStorer
		sizedIndexes  map[int]*indexShard
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &memoryIndex{
				maxKSize:      tt.fields.maxKSize,
				attributeMeta: tt.fields.attributeMeta,
				sizedIndexes:  tt.fields.sizedIndexes,
			}
			if err := k.Build(); (err != nil) != tt.wantErr {
				t.Errorf("memoryIndex.Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_memoryIndex_MaxKSize(t *testing.T) {
	type fields struct {
		maxKSize      int
		attributeMeta expr.AttributeMetadataStorer
		sizedIndexes  map[int]*indexShard
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &memoryIndex{
				maxKSize:      tt.fields.maxKSize,
				attributeMeta: tt.fields.attributeMeta,
				sizedIndexes:  tt.fields.sizedIndexes,
			}
			if got := k.MaxKSize(); got != tt.want {
				t.Errorf("memoryIndex.MaxKSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryIndex_GetPostingLists(t *testing.T) {
	type fields struct {
		maxKSize      int
		attributeMeta expr.AttributeMetadataStorer
		sizedIndexes  map[int]*indexShard
	}
	type args struct {
		size   int
		labels expr.Assignment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*PostingList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &memoryIndex{
				maxKSize:      tt.fields.maxKSize,
				attributeMeta: tt.fields.attributeMeta,
				sizedIndexes:  tt.fields.sizedIndexes,
			}
			if got := k.GetPostingLists(tt.args.size, tt.args.labels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memoryIndex.GetPostingLists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryIndex_Match(t *testing.T) {
	type fields struct {
		maxKSize      int
		attributeMeta expr.AttributeMetadataStorer
		sizedIndexes  map[int]*indexShard
	}
	type args struct {
		assignment expr.Assignment
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &memoryIndex{
				maxKSize:      tt.fields.maxKSize,
				attributeMeta: tt.fields.attributeMeta,
				sizedIndexes:  tt.fields.sizedIndexes,
			}
			if got := k.Match(tt.args.assignment); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memoryIndex.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
