package slashCommands

import "github.com/bwmarrin/discordgo"

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Sends pong along with time it took to process it.",
	},
	{
		Name:        "cat",
		Description: "Sends a random cat photo.",
	},
	{
		Name:        "help",
		Description: "Sends a list of all ",
	},
	{
		Name:        "github",
		Description: "Controls tracked GitHub repos",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "Operation",
				Description: "Operation you want to do using the GH Integration.",
				Required:    true,
			},
		},
	},
}
