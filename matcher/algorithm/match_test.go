package algorithm

import (
	"testing"
)

type createUserParams struct {
	requests     []int64
	validMatches []int64
}

func createUsers(params ...*createUserParams) []*algUser {
	users := []*algUser{}
	for i, p := range params {
		u := &algUser{
			id:               int64(i + 1),
			isMale:           true,
			requests:         p.requests,
			validMatches:     map[int64]*algUser{},
			validMatchesCopy: map[int64]*algUser{},
		}
		for _, v := range p.validMatches {
			u.validMatches[v] = nil
			u.validMatchesCopy[v] = nil
		}
		users = append(users, u)
	}
	return users
}

func TestMatch(t *testing.T) {
	matches, err := match(
		createUsers(
			&createUserParams{[]int64{5}, []int64{6}},
			&createUserParams{[]int64{5}, []int64{5, 6}},
			&createUserParams{[]int64{5}, []int64{6}},
			&createUserParams{[]int64{5}, []int64{5, 6}},
		),
		2,
		2,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 4 {
		t.Fatalf("%d != 4", len(matches))
	}
}

func TestFail(t *testing.T) {
	_, err := match(
		createUsers(
			&createUserParams{[]int64{}, []int64{3}},
			&createUserParams{[]int64{}, []int64{3}},
		),
		2,
		2,
	)
	if err == nil {
		t.Fatal("error expected")
	}
}
