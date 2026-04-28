package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "help",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Locale("ja"):    "ヘルプ",
				discordgo.Locale("en-US"): "help",
			},
			Description: "Help Command",
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Locale("ja"):    "ヘルプコマンド",
				discordgo.Locale("en-US"): "Help Command",
			},
		},
		{
			Name: "servers",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Locale("ja"):    "サーバー一覧",
				discordgo.Locale("en-US"): "servers",
			},
			Description: "List serverlist",
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Locale("ja"):    "サーバー一覧表示",
				discordgo.Locale("en-US"): "List serverlist",
			},
		},
		{
			Name: "server",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Locale("ja"):    "サーバー",
				discordgo.Locale("en-US"): "server",
			},
			Description: "Manage Server",
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Locale("ja"):    "サーバーの管理",
				discordgo.Locale("en-US"): "Manage Server",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "server_id",
					Description: "ServerId(Identifier)",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionString,
				},
				{
					Name:        "id",
					Description: "Serverid(Serverlist Id)",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionInteger,
				},
				{
					Name:        "action",
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
			Name: "role",
			NameLocalizations: &map[discordgo.Locale]string{
				discordgo.Locale("ja"):    "ロール",
				discordgo.Locale("en-US"): "role",
			},
			Description: "Manage role",
			DescriptionLocalizations: &map[discordgo.Locale]string{
				discordgo.Locale("ja"):    "ロール管理",
				discordgo.Locale("en-US"): "Manage role",
			},
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "action",
					Description: "Action that you want",
					Required:    true,
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "List",
							Value: "list",
						},
						{
							Name:  "Add",
							Value: "add",
						},
						{
							Name:  "Remove",
							Value: "remove",
						},
					},
				},
				{
					Name:        "role",
					Description: "Role that you want to manage",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionRole,
				},
				{
					Name:        "server_identifier",
					Description: "Server Identifier",
					Required:    false,
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
	}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help":    HelpCommandHandler,
		"servers": ServersCommandHandler,
		"server":  ServerCommandHandler,
		"role":    RoleCommandHandler,
	}
)

func SyncCommands(s *discordgo.Session, guildID string, appID string) {
	_, err := s.ApplicationCommandBulkOverwrite(appID, guildID, commands)
	if err != nil {
		log.Panicf("コマンドの同期に失敗しました: %v", err)
	}
	log.Println("コマンドを更新しました")
}

func OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	name := i.ApplicationCommandData().Name
	handler, ok := CommandHandlers[name]
	if !ok {
		return
	}

	handler(s, i)
}
