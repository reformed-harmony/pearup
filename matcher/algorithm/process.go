package algorithm

import (
	"errors"
)

// processParams divides the list of registered users by gender; the gender
// with the most users is returned as a list in addition to the number of users
// in the other list.
func processParams(p *Params) ([]*algUser, int, error) {
	var (
		userMajMap = map[int64]*algUser{}
		userMinMap = map[int64]*algUser{}
	)

	// Populate the map for each of the genders
	for _, r := range p.Registrants {
		u := &algUser{
			id:     r.ID,
			isMale: r.IsMale,
		}
		if r.IsMale {
			userMajMap[r.ID] = u
		} else {
			userMinMap[r.ID] = u
		}
	}

	// Make sure neither map is empty
	if len(userMajMap) == 0 || len(userMinMap) == 0 {
		return nil, 0, errors.New("one or more lists is empty")
	}

	// Ensure that userMajMap is assigned the larger list
	if len(userMinMap) > len(userMajMap) {
		userMinMap, userMajMap = userMajMap, userMinMap
	}

	// For every user in userMajMap, assign a copy of the other map
	for _, u := range userMajMap {
		u.validMatches = copyMap(userMinMap)
	}

	// Apply the requests to the userMajMap; the requests by those in
	// userMinMap are inverted and assigned to userMajMap
	for _, r := range p.Requests {
		if u := userMajMap[r.User1ID]; u != nil {
			u.requests = append(u.requests, r.User2ID)
		}
		if u := userMajMap[r.User2ID]; u != nil {
			u.requests = append(u.requests, r.User1ID)
		}
	}

	// Apply the exclusions by removing the items from validMatches
	for _, e := range p.Exclusions {
		if u := userMajMap[e.User1ID]; u != nil {
			delete(u.validMatches, e.User2ID)
		}
		if u := userMajMap[e.User2ID]; u != nil {
			delete(u.validMatches, e.User1ID)
		}
	}

	// Convert userMajMap into a list and assign a copy of validMatches
	userMajList := []*algUser{}
	for _, u := range userMajMap {
		userMajList = append(userMajList, u)
		u.validMatchesCopy = copyMap(u.validMatches)
	}

	// Return the list and size of the userMinMap
	return userMajList, len(userMinMap), nil
}
