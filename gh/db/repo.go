package db

import (
	"time"
)

type Repo struct {
	Owner       string `gorm:"primaryKey"`
	Name        string `gorm:"primaryKey"`
	LastIssue   time.Time
	TrackIssues bool
}
