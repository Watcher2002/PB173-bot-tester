package gh

import (
	"PB173-discord-bot/gh/db"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/v56/github"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func CheckForIssues() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		tracked, err := db.GetIssueTrackedRepos(RepoDB)
		if err != nil {
			log.Error().Msg("Couldn't get Issue Tracked Repos.")
		}

		for _, repo := range tracked {
			go getNewIssues(repo)
		}
		<-ticker.C
	}
}

func GetIssues(_ *github.Client, owner, repo string) []*github.Issue {
	issues, _, err := Client.Issues.ListByRepo(context.Background(), owner, repo, nil)
	if err != nil {
		log.Error().Msg("Couldn't get issues from repo.")
	}

	return issues
}

func postNewIssue(issue *github.Issue) {
	fmt.Print(issue.GetRepository())

	title := fmt.Sprintf("New issue in %s", issue.GetRepository().GetFullName())
	desc := fmt.Sprintf("[%s](%s)", issue.GetTitle(), issue.GetHTMLURL())
	channelID := os.Getenv("GITHUB_CHANNEL")
	if channelID == "" {
		log.Error().Msg("Missing github channel ID.")
		return
	}

	Session.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
	})
}

func getNewIssues(repo *db.Repo) {
	issues, _, err := Client.Issues.ListByRepo(context.Background(), repo.Owner, repo.Name, &github.IssueListByRepoOptions{
		Since: repo.LastIssue.Add(time.Second),
	})

	if err != nil {
		msg := fmt.Sprintf("Failed to get %s/%s issues", repo.Owner, repo.Name)
		log.Error().Msg(msg)
	}

	if len(issues) == 0 {
		return
	}

	for _, issue := range issues {
		postNewIssue(issue)
	}

	err = db.UpdateLastIssue(RepoDB, repo, issues[len(issues)-1].CreatedAt.Time)
	if err != nil {
		log.Error().Msg("Couldn't update " + repo.Owner + "/" + repo.Name)
	}

}
