package gh

import (
	"PB173-discord-bot/gh/db"
	"context"
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/v56/github"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	Client  *github.Client
	Session *discordgo.Session
	RepoDB  *gorm.DB
)

func ConnectToGithub() error {
	Client = github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	if Client == nil {
		log.Fatal().Msg("Couldn't connect to GitHub.")
		return errors.New("couldn't connect to GitHub")
	}

	RepoDB = db.ConnDB()
	log.Info().Msg("Connected to database.")
	return nil
}

func GetRepo(URL string) *db.Repo {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return nil
	}

	segments := strings.Split(parsedURL.Path, "/")
	if len(segments) < 3 {
		return nil
	}

	owner := segments[1]
	repo := segments[2]

	_, _, err = Client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return nil
	}

	return &db.Repo{
		Owner:       owner,
		Name:        repo,
		TrackIssues: true,
		LastIssue:   time.Now(),
	}
}
