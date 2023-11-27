package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func EmojiReactionHandler(dc *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	if reaction.UserID == dc.State.User.ID {
		return
	}

	switch reaction.Emoji.Name {
	case "🔖":
		channel, err := dc.UserChannelCreate(reaction.UserID)
		if err != nil {
			log.Error().Msg("Couldn't connect to a Users' channel.")
			return
		}
		link := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", reaction.GuildID, reaction.ChannelID, reaction.MessageID)
		_, err = dc.ChannelMessageSend(channel.ID, "You have bookmarked this post:\n"+link)
		if err != nil {
			log.Error().Msg("Couldn't send user bookmark.")
		}
	default:
		return
	}
}
