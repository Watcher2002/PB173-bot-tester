package events

import "github.com/bwmarrin/discordgo"

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.User.Bot || i.Data.Type() != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ping":
		go SlashPing(s, i)
	}
}
