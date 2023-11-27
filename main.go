package main

import (
	"PB173-discord-bot/bot"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error().Msg("Couldn't load environment.")
	} else {
		log.Info().Msg("Environment var loading successful.")
	}

	bot.StartBot()
}
