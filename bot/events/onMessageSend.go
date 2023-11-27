package events

import (
	"PB173-discord-bot/gh"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

func MessageHandler(dc *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == dc.State.User.ID {
		return
	}

	prefix := os.Getenv("PREFIX")
	if prefix == "" {
		log.Info().Msg("Missing bot prefix")
		return
	}

	if !strings.HasPrefix(message.Content, prefix) {
		return
	}
	if strings.HasPrefix(strings.ToLower(message.Content[len(prefix):]), "github") {
		go gh.GithubArgParser(dc, message)
		return
	}

	if strings.HasPrefix(strings.ToLower(message.Content[len(prefix):]), "wiki") {
		go wikiSummary(dc, message)
	}

	switch strings.ToLower(message.Content[len(prefix):]) {
	case "help":
		go help(dc, message)
	case "ping":
		go ping(dc, message)
	case "cat":
		go getRandomCat(dc, message)
	}
}
