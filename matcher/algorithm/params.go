package algorithm

// Request represents a match that the first user would like to take place.
type Request struct {
	User1ID int64
	User2ID int64
}

// Exclusion represents a match that cannot take place. This could be caused by
// a previous match or an explicit request from a user.
type Exclusion struct {
	User1ID int64
	User2ID int64
}

// Registrant represents a user that has registered for being matched.
type Registrant struct {
	ID     int64
	IsMale bool
}

// Params represents the input data for the algorithm.
type Params struct {
	Requests    []*Request
	Exclusions  []*Exclusion
	Registrants []*Registrant
}

// Match represents a match determined by the algorithm.
type Match struct {
	User1ID int64
	User2ID int64
}
