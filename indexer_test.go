package main

import (
	"testing"
)

func Test_kIndexTable_Add(t *testing.T) {
	k := newKIndexTable()

	k.Add(NewConjunction(
		1,
		[]*Attribute{
			{Key: "age", Values: []string{"3"}, Contains: true},
			{Key: "state", Values: []string{"NY"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		2,
		[]*Attribute{
			{Key: "age", Values: []string{"3"}, Contains: true},
			{Key: "gender", Values: []string{"F"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		3,
		[]*Attribute{
			{Key: "age", Values: []string{"3"}, Contains: true},
			{Key: "gender", Values: []string{"M"}, Contains: true},
			{Key: "state", Values: []string{"CA"}, Contains: false},
		},
	))

	k.Add(NewConjunction(
		4,
		[]*Attribute{
			{Key: "state", Values: []string{"CA"}, Contains: true},
			{Key: "gender", Values: []string{"M"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		5,
		[]*Attribute{
			{Key: "age", Values: []string{"3", "4"}, Contains: true},
		},
	))

	k.Add(NewConjunction(
		6,
		[]*Attribute{
			{Key: "state", Values: []string{"CA", "NY"}, Contains: false},
		},
	))

	k.Build()

	k.MaxKSize()

}

func Test_sortPostingList(t *testing.T) {
	p := PostingList{
		&PostingItem{ConjunctionID: 2, Contains: true},
		&PostingItem{ConjunctionID: 1, Contains: true},
		&PostingItem{ConjunctionID: 3, Contains: true},
		&PostingItem{ConjunctionID: 1, Contains: false},
	}

	sortPostingList(p)

}
