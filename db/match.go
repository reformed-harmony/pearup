package db

// Match represents a pairing between two users.
type Match struct {
	ID       int64
	Pearup   *Pearup `gorm:"ForeignKey:PearupID"`
	PearupID int64   `sql:"type:int REFERENCES pearups(id) ON DELETE CASCADE"`
	User1    *User   `gorm:"ForeignKey:User1ID"`
	User1ID  int64   `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
	User2    *User   `gorm:"ForeignKey:User2ID"`
	User2ID  int64   `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
}
