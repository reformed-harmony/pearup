package db

// Registration represents an intent by a user to participate in a pearup.
type Registration struct {
	ID       int64
	Pearup   *Pearup `gorm:"ForeignKey:PearupID"`
	PearupID int64   `sql:"type:int REFERENCES pearups(id) ON DELETE CASCADE"`
	User     *User   `gorm:"ForeignKey:UserID"`
	UserID   int64   `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
}
