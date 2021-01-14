package db

// Request represents a desire by a user to be matched with a specific user.
type Request struct {
	ID              int64
	User            *User `gorm:"ForeignKey:UserID"`
	UserID          int64 `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
	RequestedUser   *User `gorm:"ForeignKey:RequestedUserID"`
	RequestedUserID int64 `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
}
