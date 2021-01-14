package algorithm

import (
	"reflect"
	"testing"
)

func TestSimple(t *testing.T) {
	userList, minTotal, err := processParams(
		&Params{
			Registrants: []*Registrant{
				{1, true},
				{2, true},
				{3, false},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(userList) != 2 {
		t.Fatalf("%d != 2", len(userList))
	}
	if minTotal != 1 {
		t.Fatalf("%d != 1", minTotal)
	}
	for i, userID := range []int{1, 2} {
		u := userList[i]
		if len(u.validMatches) != 1 {
			t.Fatalf("#%d: %d != 1", userID, len(u.validMatches))
		}
		for _, m := range u.validMatches {
			if m.id != 3 {
				t.Fatalf("#%d: %d != 3", userID, m.id)
			}
		}
		if !reflect.DeepEqual(u.validMatches, u.validMatchesCopy) {
			t.Fatalf("#%d: %v != %v", userID, u.validMatches, u.validMatchesCopy)
		}
	}
}

func TestExlusion(t *testing.T) {
	userList, _, err := processParams(
		&Params{
			Exclusions: []*Exclusion{
				{2, 3},
			},
			Registrants: []*Registrant{
				{1, true},
				{2, true},
				{3, false},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(userList) != 2 {
		t.Fatalf("%d != 2", len(userList))
	}
	if len(userList[1].validMatches) != 0 {
		t.Fatalf("%d != 1", len(userList[1].validMatches))
	}
}
