package algorithm

import (
	"errors"
)

var errNoMatches = errors.New("no valid matches")

// match is a recursive function (yes, I know)
func match(majList []*algUser, minRemaining, minTotal int) ([]*Match, error) {

	var (
		u            = majList[0]
		validMatches []int64
	)

	// Add valid requests to validMatches
	for _, r := range u.requests {
		if _, ok := u.validMatches[r]; ok {
			validMatches = append(validMatches, r)
			delete(u.validMatches, r)
		}
	}

	// Add what remains in u.validMatches
	for v := range u.validMatches {
		validMatches = append(validMatches, v)
	}

	// Loop over valid matches
	for _, userID := range validMatches {

		// Create the match
		m := &Match{
			User1ID: u.id,
			User2ID: userID,
		}

		// If this is the last match, return immediately
		if len(majList) == 1 {
			return []*Match{m}, nil
		}

		// Remove userID from validMatches in majList
		var newMajList []*algUser
		for _, u := range majList[1:] {

			// Build a new list of validMatches for this user
			var newValidMatches map[int64]*algUser
			if minRemaining == 1 {
				newValidMatches = copyMap(u.validMatchesCopy)
			} else {
				newValidMatches = copyMap(u.validMatches)
				delete(newValidMatches, userID)
			}

			// Create the new algUser struct and add it to the list
			newUser := &algUser{
				id:               u.id,
				isMale:           u.isMale,
				requests:         u.requests,
				validMatches:     newValidMatches,
				validMatchesCopy: u.validMatchesCopy,
			}
			newMajList = append(newMajList, newUser)
		}

		// If we have completed a cycle through minList, reset it
		if minRemaining == 0 {
			minRemaining = minTotal
		} else {
			minRemaining--
		}

		// Continue making matches
		matches, err := match(newMajList, minRemaining, minTotal)
		if err != nil {
			continue
		}

		// Everything down the call chain matched; return the matches
		return append(matches, m), nil
	}

	// There were no valid matches
	return nil, errNoMatches
}
