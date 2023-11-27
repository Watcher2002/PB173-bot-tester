package events

import (
	"PB173-discord-bot/bot/ELI5"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func help(dc *discordgo.Session, message *discordgo.MessageCreate) {
	_, err := dc.ChannelMessageSendEmbedReply(message.ChannelID, &discordgo.MessageEmbed{
		Title:       "Help",
		Description: "Current prefix: **" + os.Getenv("PREFIX") + "**\nI can help you with few things:\n- **ping** - How long it took for the message to arrive\n- **cat** - Sends a picture of a random cat\n- **github** [add/remove/track/list] For add and remove give link, for track give what to change[issue], \"on/off\" and link\n- **wiki** - Search for basic things such as cat or car :)",
	}, message.Reference())
	if err != nil {
		log.Error().Msg("Message couldn't be sent.")
	}
}

func ping(dc *discordgo.Session, message *discordgo.MessageCreate) {
	desc := fmt.Sprintf("@%s (%d ms)", message.Author.GlobalName, discordEpochToUNIX(message.ID))
	_, err := dc.ChannelMessageSendEmbedReply(message.ChannelID, &discordgo.MessageEmbed{
		Title:       "Pong",
		Description: desc,
		Author: &discordgo.MessageEmbedAuthor{
			URL:          "",
			Name:         dc.State.User.Username,
			IconURL:      "",
			ProxyIconURL: "",
		},
	}, message.Reference())
	if err != nil {
		log.Error().Msg("Message couldn't be sent.")
	}
}

func discordEpochToUNIX(ID string) uint64 {
	id, _ := strconv.ParseUint(ID, 10, 64)
	return (uint64)(time.Now().UnixMilli()) - ((id >> 22) + 1420070400000)
}

const catURL = "https://cataas.com/"

func getRandomCat(dc *discordgo.Session, message *discordgo.MessageCreate) {
	resp, err := http.Get(catURL + "cat")
	if err != nil {
		log.Error().Msg("Couldn't get cat :(")
		dc.ChannelMessageSendReply(message.ChannelID, "Sadly there was an error getting your cat :(", message.Reference())
		return
	}
	defer resp.Body.Close()

	dc.ChannelFileSend(message.ChannelID, "Kitty for you.jpg", resp.Body)
}

func wikiSummary(dc *discordgo.Session, message *discordgo.MessageCreate) {
	args := strings.Split(message.Content, " ")

	if len(args) == 1 {
		dc.ChannelMessageSendReply(message.ChannelID, "What do you want to search for?", message.Reference())
		return
	}

	info, err := ELI5.GetWikiArticleExtract(strings.Join(args[1:], " "))
	if err != nil {
		return
	}

	dc.ChannelMessageSendReply(message.ChannelID, info, message.Reference())

}
