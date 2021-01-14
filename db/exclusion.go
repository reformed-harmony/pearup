package db

// Exclusion represents a desire by a user to avoid being paired with another user.
type Exclusion struct {
	ID             int64
	User           *User `gorm:"ForeignKey:UserID"`
	UserID         int64 `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
	ExcludedUser   *User `gorm:"ForeignKey:ExcludedUserID"`
	ExcludedUserID int64 `sql:"type:int REFERENCES users(id) ON DELETE CASCADE"`
}
