package discord

import (
	""
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			name: "Help",
			Description: "Help Command"
		},
		{
			name: "Servers",
			Description: "List serverlist"
		},
		{
			name: "Server",
			Description: "Manage Server"
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name: "ServerId"
					Description: "ServerId(Identifier)"
					Required: False,
					Type: discordgo.ApplicationCommandOptionString,
				},
				{
					Name: "Id"
					Description: "Serverid(Serverlist Id)"
					Required: False
					Type: discordgo.ApplicationCommandOptionInteger
				},
				{
					Name: "Action"
					Description "Action that you want"
					Required: False
					Type: discordgo.ApplicationCommandOptionString
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Start"
							Description: "Start server"
							Value: "start"
						},
						{
							Name: "Stop"
							Description: "Stop server"
							Value: "stop"
						},
						{
							Name: "Restart"
							Description: "Restart server"
							Value: "restart"
						},
						{
							Name: "Information"
							Description: "Server specific"
							Value: "info"
						},
					},
				},
			},

		},
		{
			name: "Roll",
			Description: "Manage roll"
			
		}
	}
)
