package db

import (
	"time"
)

// Pearup represents a single pearup.
type Pearup struct {
	ID         int64
	Name       string    `gorm:"not null"`
	EndDate    time.Time `gorm:"not null"`
	IsComplete bool      `gorm:"not null"`
	IsPublic   bool      `gorm:"not null"`
	CanRequest bool      `gorm:"not null" sql:"DEFAULT:FALSE"`
}
