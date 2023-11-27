package gh

import (
	"PB173-discord-bot/gh/db"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"strings"
)

func GithubArgParser(dc *discordgo.Session, message *discordgo.MessageCreate) {
	if Client == nil {
		_, err := dc.ChannelMessageSendReply(message.ChannelID, "gh integration has not been configured.", message.Reference())
		if err != nil {
			log.Error().Msg("Message couldn't be sent.")
		}
	}

	arguments := strings.Split(strings.ToLower(message.Content), " ")
	if len(arguments) == 1 {
		unknownCommandReply(dc, message)
		return
	}

	switch arguments[1] {
	case "add":
		addGithubRepos(dc, message, arguments)
	case "remove":
		removeGithubRepos(dc, message, arguments)
	case "list":
		listTrackedRepos(dc, message, arguments)
	case "track":
		changeTracking(dc, message, arguments)
	default:
		unknownCommandReply(dc, message)
	}
}

func changeTracking(dc *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	if len(args) != 5 {
		unknownCommandReply(dc, message)
		return
	}

	opt := args[3]
	repo := GetRepo(args[4])

	if repo == nil {
		log.Info().Msg("Unknown repo: " + args[4])
		dc.ChannelMessageSendReply(message.ChannelID, "Repo could't be found.", message.Reference())
		return
	}

	switch args[2] {
	case "issue":
		switch opt {
		case "on":
			repo.TrackIssues = true
		case "off":
			repo.TrackIssues = false
		default:
			return
		}
	}

	res := RepoDB.Save(repo)
	if res.Error != nil {
		log.Error().Msg("Couldn't update tracking of repo.")
	}
	err := dc.MessageReactionAdd(message.ChannelID, message.ID, "✔️")
	if err != nil {
		log.Error().Msg("Emoji couldn't be sent.")
		return
	}

}

func listTrackedRepos(dc *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	var toPrint []*db.Repo
	var err error
	if len(args) == 2 {
		toPrint, err = db.GetAllRepos(RepoDB)

		if err != nil {
			log.Error().Msg("Couldn't get all repos.")
			return
		}

	} else if len(args) == 3 {
		if args[2] != "tracked" {
			unknownCommandReply(dc, message)
			return
		}
		toPrint, err = db.GetIssueTrackedRepos(RepoDB)

		if err != nil {
			log.Error().Msg("Couldn't get tracked repos.")
			return
		}

	} else {
		unknownCommandReply(dc, message)
		return
	}

	var desc string
	for _, repo := range toPrint {
		desc += fmt.Sprintf("Owner: %s\tName:%s\n", repo.Owner, repo.Name)
	}

	if desc == "" {
		desc = "There are no tracked repos, you can add them by using !github add"
	}

	dc.ChannelMessageSendEmbedReply(message.ChannelID, createEmbed("Tracked repos:", desc), message.Reference())
}

func removeGithubRepos(dc *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		unknownCommandReply(dc, message)
		return
	}

	for i := 2; i < len(args); i++ {
		repo := GetRepo(args[i])
		if db.RemoveRepo(RepoDB, repo) != nil {
			desc := fmt.Sprintf("Was %s tracked before?", args[i])
			dc.ChannelMessageSendEmbedReply(message.ChannelID, createEmbed("Failed to remove repo", desc), message.Reference())
			log.Error().Msg("Couldn't remove repo from DB.")
		} else {
			log.Info().Msg("Removed repo from DB")
		}
	}

	dc.MessageReactionAdd(message.ChannelID, message.ID, "✔️")
}

func addGithubRepos(dc *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		unknownCommandReply(dc, message)
		return
	}

	suc, fail := addReposToDB(args[2:])
	for _, failedLinks := range fail {
		desc := fmt.Sprintf("The link (%s) doesn't seem correct.", failedLinks)
		dc.ChannelMessageSendEmbedReply(message.ChannelID, createEmbed("Couldn't add repo", desc), message.Reference())
		msg := fmt.Sprintf("Couldn't add %s", failedLinks)
		log.Error().Msg(msg)
	}

	for _, successful := range suc {
		msg := fmt.Sprintf("Added %s", successful)
		log.Info().Msg(msg)
	}

	if len(fail) == 0 {
		dc.MessageReactionAdd(message.ChannelID, message.ID, "✔️")
	}
}

func addReposToDB(links []string) (success, failed []string) {
	failed = []string{}
	success = []string{}
	fmt.Printf("%v", links)
	for _, link := range links {
		repo := GetRepo(link)
		if repo == nil {
			failed = append(failed, link)
			continue
		}
		db.AddRepo(RepoDB, repo)
		success = append(success, link)
	}

	return success, failed
}

func unknownCommandReply(dc *discordgo.Session, m *discordgo.MessageCreate) {
	dc.ChannelMessageSendEmbedReply(m.ChannelID, createEmbed("Unknown command", "See !help for more info."), m.Reference())
}

func createEmbed(title, desc string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Author: &discordgo.MessageEmbedAuthor{
			Name: Session.State.User.Username,
		},
	}
}
