package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_kIndexTable_Add(t *testing.T) {
	k := newKIndexTable()

	k.Add(NewConjunction(
		1,
		[]*Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true},
			{Name: "state", Values: []string{"NY"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		2,
		[]*Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true},
			{Name: "gender", Values: []string{"F"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		3,
		[]*Attribute{
			{Name: "age", Values: []string{"3"}, Contains: true},
			{Name: "gender", Values: []string{"M"}, Contains: true},
			{Name: "state", Values: []string{"CA"}, Contains: false},
		},
	))

	k.Add(NewConjunction(
		4,
		[]*Attribute{
			{Name: "state", Values: []string{"CA"}, Contains: true},
			{Name: "gender", Values: []string{"M"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		5,
		[]*Attribute{
			{Name: "age", Values: []string{"3", "4"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		6,
		[]*Attribute{
			{Name: "state", Values: []string{"CA", "NY"}, Contains: false},
		},
	))

	k.Build()

	assert.Equal(t, 2, k.MaxKSize())
	zeroIdx := k.sizedIndexes[0].(*memoryIndexer)
	assert.Equal(t, 3, len(zeroIdx.imap))
	assert.Equal(t, int64(6), zeroIdx.imap["state:CA"].Items[0].CID)
	assert.Equal(t, false, zeroIdx.imap["state:CA"].Items[0].Contains)
	assert.Equal(t, int64(6), zeroIdx.imap["state:NY"].Items[0].CID)
	assert.Equal(t, false, zeroIdx.imap["state:NY"].Items[0].Contains)
	assert.Equal(t, int64(6), zeroIdx.imap["null:null"].Items[0].CID)
	assert.Equal(t, true, zeroIdx.imap["null:null"].Items[0].Contains)

	oneIdx := k.sizedIndexes[1].(*memoryIndexer)
	assert.Equal(t, 2, len(oneIdx.imap))
	assert.Equal(t, int64(5), oneIdx.imap["age:3"].Items[0].CID)
	assert.Equal(t, true, oneIdx.imap["age:3"].Items[0].Contains)
	assert.Equal(t, int64(5), oneIdx.imap["age:4"].Items[0].CID)
	assert.Equal(t, true, oneIdx.imap["age:4"].Items[0].Contains)

	twoIdx := k.sizedIndexes[2].(*memoryIndexer)
	assert.Equal(t, 5, len(twoIdx.imap))
	assert.Equal(t, int64(1), twoIdx.imap["state:NY"].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap["state:NY"].Items[0].Contains)

	assert.Equal(t, int64(1), twoIdx.imap["age:3"].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap["age:3"].Items[0].Contains)
	assert.Equal(t, int64(2), twoIdx.imap["age:3"].Items[1].CID)
	assert.Equal(t, true, twoIdx.imap["age:3"].Items[1].Contains)
	assert.Equal(t, int64(3), twoIdx.imap["age:3"].Items[2].CID)
	assert.Equal(t, true, twoIdx.imap["age:3"].Items[2].Contains)

	assert.Equal(t, int64(2), twoIdx.imap["gender:F"].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap["gender:F"].Items[0].Contains)

	assert.Equal(t, int64(3), twoIdx.imap["state:CA"].Items[0].CID)
	assert.Equal(t, false, twoIdx.imap["state:CA"].Items[0].Contains)
	assert.Equal(t, int64(4), twoIdx.imap["state:CA"].Items[1].CID)
	assert.Equal(t, true, twoIdx.imap["state:CA"].Items[1].Contains)

	assert.Equal(t, int64(3), twoIdx.imap["gender:M"].Items[0].CID)
	assert.Equal(t, true, twoIdx.imap["gender:M"].Items[0].Contains)
	assert.Equal(t, int64(4), twoIdx.imap["gender:M"].Items[1].CID)
	assert.Equal(t, true, twoIdx.imap["gender:M"].Items[1].Contains)

	matcher := &matcher{}
	matched := matcher.Match(k, Labels{
		Label{Name: "age", Value: "3"},
		Label{Name: "state", Value: "CA"},
		Label{Name: "gender", Value: "M"},
	})

	assert.Equal(t, []int64{int64(4), int64(5)}, matched)
}
