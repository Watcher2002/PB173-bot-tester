package bot

import (
	"PB173-discord-bot/bot/events"
	"PB173-discord-bot/gh"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

var (
	Dc *discordgo.Session
)

func StartBot() {
	var err error
	Dc, err = discordgo.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal().Msg("Couldn't create session.")
	}

	Dc.SyncEvents = false
	Dc.AddHandler(events.MessageHandler)
	Dc.AddHandler(events.EmojiReactionHandler)
	Dc.AddHandler(events.InteractionHandler)

	err = Dc.Open()

	Dc.ApplicationCommandCreate(Dc.State.User.ID, "1175785400611131409", &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Basic Command",
	})

	defer Dc.Close()
	if err != nil {
		log.Fatal().Msg("Couldn't establish connection.")
	}
	log.Info().Msg("Bot has started.")

	err = gh.ConnectToGithub()
	if err != nil {
		return
	}
	gh.Session = Dc

	go gh.CheckForIssues()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
