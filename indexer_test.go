package main

import (
	"testing"
)

func Test_kIndexTable_Add(t *testing.T) {
	k := newKIndexTable()

	k.Add(NewConjunction(
		1,
		[]*Attribute{
			{Key: "age", Value: []string{"3"}, BelongsTo: true},
			{Key: "state", Value: []string{"NY"}, BelongsTo: true},
		},
	))

	k.Add(NewConjunction(
		2,
		[]*Attribute{
			{Key: "age", Value: []string{"3"}, BelongsTo: true},
			{Key: "gender", Value: []string{"F"}, BelongsTo: true},
		},
	))

	k.Add(NewConjunction(
		3,
		[]*Attribute{
			{Key: "age", Value: []string{"3"}, BelongsTo: true},
			{Key: "gender", Value: []string{"M"}, BelongsTo: true},
			{Key: "state", Value: []string{"CA"}, BelongsTo: false},
		},
	))

	k.Add(NewConjunction(
		4,
		[]*Attribute{
			{Key: "state", Value: []string{"CA"}, BelongsTo: true},
			{Key: "gender", Value: []string{"M"}, BelongsTo: true},
		},
	))

	k.Add(NewConjunction(
		5,
		[]*Attribute{
			{Key: "age", Value: []string{"3", "4"}, BelongsTo: true},
		},
	))

	k.Add(NewConjunction(
		6,
		[]*Attribute{
			{Key: "state", Value: []string{"CA", "NY"}, BelongsTo: false},
		},
	))

	k.Build()

	k.MaxKSize()

}
