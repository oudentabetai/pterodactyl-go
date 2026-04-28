package discord

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/oudentabetai/pterodactyl-go/utils"
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
					Name:         "server_name",
					Description:  "Server Name",
					Required:     false,
					Autocomplete: true,
					Type:         discordgo.ApplicationCommandOptionString,
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
	if i.Type == discordgo.InteractionApplicationCommand {
		name := i.ApplicationCommandData().Name
		handler, ok := CommandHandlers[name]
		if !ok {
			return
		}
		handler(s, i)
	} else {
		if i.Type == discordgo.InteractionApplicationCommandAutocomplete {
			data := i.ApplicationCommandData()

			var choices []*discordgo.ApplicationCommandOptionChoice

			for _, opt := range data.Options {
				if opt.Focused {
					userInput := opt.StringValue()

					servers := utils.GetAccessibleServers(i.Member)
					lowerInput := strings.ToLower(userInput)
					for _, srv := range servers {
						name := srv.Attributes.Name
						if lowerInput == "" || strings.Contains(strings.ToLower(name), lowerInput) {
							choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
								Name:  name,
								Value: srv.Attributes.Identifier,
							})
							if len(choices) >= 25 {
								break
							}
						}
					}

					if len(choices) > 0 {
						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseType(8), // Application Command Autocomplete Result
							Data: &discordgo.InteractionResponseData{
								Choices: choices,
							},
						})
						if err != nil {
							log.Printf("autocomplete respond error: %v", err)
						}
					}

					return
				}
			}
		}
	}

}
