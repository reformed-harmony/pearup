package algorithm

import (
	"errors"
	"sort"
)

var (
	errMatchMismatch = errors.New("number of matches does not equal expected value")
)

// Algorithm maintains state during the matching process.
type Algorithm struct {
	params *Params
}

// New creates a new Algorithm instance.
func New(params *Params) *Algorithm {
	return &Algorithm{
		params: params,
	}
}

// algUser maintains information for a specific user.
type algUser struct {
	id               int64
	isMale           bool
	requests         []int64
	validMatches     map[int64]*algUser
	validMatchesCopy map[int64]*algUser
}

// byRequests sorts algUsers by whether they have requests or not
type byRequests []*algUser

func (r byRequests) Len() int           { return len(r) }
func (r byRequests) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r byRequests) Less(i, j int) bool { return len(r[i].requests) >= 0 && len(r[j].requests) == 0 }

// Run executes the algorithm and attempts to produce a list of matches. The
// algorithm is an attempt to take all possible valid "solutions" to the input
// parameters and find the first valid solution as quickly as possible. A stack
// is used for "pushing" and "pulling" solutions in progress.
func (a *Algorithm) Run() ([]*Match, error) {

	// Process the parameters (split into list and map)
	majList, minTotal, err := processParams(a.params)
	if err != nil {
		return nil, err
	}

	// Sort the majList by requests
	sort.Sort(byRequests(majList))

	// Perform the matches
	m, err := match(majList, minTotal, minTotal)
	if err != nil {
		return nil, err
	}

	// Perform a last minute sanity check
	if len(majList) != len(m) {
		return nil, errMatchMismatch
	}

	return m, nil
}
