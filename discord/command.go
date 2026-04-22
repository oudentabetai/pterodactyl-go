package discord

import "github.com/bwmarrin/discordgo"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "Help",
			Description: "Help Command",
		},
		{
			Name:        "Servers",
			Description: "List serverlist",
		},
		{
			Name:        "Server",
			Description: "Manage Server",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "ServerId",
					Description: "ServerId(Identifier)",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionString,
				},
				{
					Name:        "Id",
					Description: "Serverid(Serverlist Id)",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionInteger,
				},
				{
					Name:        "Action",
					Description: "Action that you want",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Start",
							Value: "start",
						},
						{
							Name:  "Stop",
							Value: "stop",
						},
						{
							Name:  "Restart",
							Value: "restart",
						},
						{
							Name:  "Information",
							Value: "info",
						},
					},
				},
			},
		},
		{
			Name:        "Roll",
			Description: "Manage roll",
		},
	}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help":    HelpCommandHandler,
		"servers": ServersCommandHandler,
		"server":  ServerCommandHandler,
		"role":    RoleCommandHandler,
	}
)
