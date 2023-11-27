package db

import (
	"github.com/glebarez/sqlite"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

func ConnDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("repos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal().Msg("Couldn't open db.")
		return nil
	}

	err = db.AutoMigrate(&Repo{})
	if err != nil {
		log.Fatal().Msg("Couln't migrate.")
	}

	return db
}

func AddRepo(db *gorm.DB, repo *Repo) error {
	res := db.Create(repo)
	return res.Error
}

func GetIssueTrackedRepos(db *gorm.DB) ([]*Repo, error) {
	var repos []*Repo
	res := db.Where("track_issues == true").Find(&repos)

	return repos, res.Error
}

func RemoveRepo(db *gorm.DB, repo *Repo) error {
	res := db.Delete(repo)
	return res.Error
}

func GetAllRepos(db *gorm.DB) ([]*Repo, error) {
	var repos []*Repo
	res := db.Find(&repos)

	return repos, res.Error
}

func UpdateLastIssue(db *gorm.DB, repo *Repo, timestamp time.Time) error {
	repo.LastIssue = timestamp
	res := db.Save(repo)

	return res.Error
}
