package matcher

import (
	"github.com/reformed-harmony/pearup/db"
	"github.com/reformed-harmony/pearup/matcher/algorithm"
)

func loadRequests(conn *db.Conn) ([]*algorithm.Request, error) {
	var (
		dbRequests = []*db.Request{}
		requests   = []*algorithm.Request{}
	)
	if err := conn.
		Find(&dbRequests).
		Error; err != nil {
		return nil, err
	}
	for _, r := range dbRequests {
		requests = append(requests, &algorithm.Request{
			User1ID: r.UserID,
			User2ID: r.RequestedUserID,
		})
	}
	return requests, nil
}

func loadExclusions(conn *db.Conn) ([]*algorithm.Exclusion, error) {
	var (
		dbMatches    = []*db.Match{}
		dbExclusions = []*db.Exclusion{}
		exclusions   = []*algorithm.Exclusion{}
	)
	if err := conn.
		Find(&dbMatches).
		Error; err != nil {
		return nil, err
	}
	for _, m := range dbMatches {
		exclusions = append(exclusions, &algorithm.Exclusion{
			User1ID: m.User1ID,
			User2ID: m.User2ID,
		})
	}
	if err := conn.
		Find(&dbExclusions).
		Error; err != nil {
		return nil, err
	}
	for _, e := range dbExclusions {
		exclusions = append(exclusions, &algorithm.Exclusion{
			User1ID: e.UserID,
			User2ID: e.ExcludedUserID,
		})
	}
	return exclusions, nil
}

func loadRegistrants(conn *db.Conn, pearupID int64) ([]*algorithm.Registrant, error) {
	var (
		dbRegistration = []*db.Registration{}
		registrants    = []*algorithm.Registrant{}
	)
	if err := conn.
		Where("pearup_id = ?", pearupID).
		Preload("User").
		Find(&dbRegistration).Error; err != nil {
		return nil, err
	}
	for _, r := range dbRegistration {
		registrants = append(registrants, &algorithm.Registrant{
			ID:     r.User.ID,
			IsMale: r.User.IsMale,
		})
	}
	return registrants, nil
}

func createMatches(conn *db.Conn, pearupID int64, matches []*algorithm.Match) error {
	for _, m := range matches {
		if err := conn.Save(&db.Match{
			PearupID: pearupID,
			User1ID:  m.User1ID,
			User2ID:  m.User2ID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func removeRequests(conn *db.Conn, matches []*algorithm.Match) error {
	dbRequests := []*db.Request{}
	if err := conn.
		Find(&dbRequests).
		Error; err != nil {
		return err
	}
	for _, m := range matches {
		for _, r := range dbRequests {
			if r.UserID == m.User1ID && r.RequestedUserID == m.User2ID ||
				r.UserID == m.User2ID && r.RequestedUserID == m.User1ID {
				if err := conn.Delete(r).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *Matcher) match(conn *db.Conn, p *db.Pearup) error {
	m.log.Infof("begining %s...", p.Name)

	// Load requests (if allowed)
	requests, err := func() ([]*algorithm.Request, error) {
		if p.CanRequest {
			return loadRequests(conn)
		}
		return []*algorithm.Request{}, nil
	}()
	if err != nil {
		return err
	}

	// Load exclusions
	exclusions, err := loadExclusions(conn)
	if err != nil {
		return err
	}

	// Load registrants
	registrants, err := loadRegistrants(conn, p.ID)
	if err != nil {
		return err
	}

	// Create the algorithm object and run it
	matches, err := algorithm.New(&algorithm.Params{
		Requests:    requests,
		Exclusions:  exclusions,
		Registrants: registrants,
	}).Run()
	if err != nil {
		return err
	}

	// Create the matches
	if err := createMatches(conn, p.ID, matches); err != nil {
		return err
	}

	// Remove requests that were completed
	if err := removeRequests(conn, matches); err != nil {
		return err
	}

	// Indicate the pear-up is complete
	p.IsComplete = true
	if err := conn.Save(p).Error; err != nil {
		return err
	}

	m.log.Infof("completed %s", p.Name)
	return nil
}
