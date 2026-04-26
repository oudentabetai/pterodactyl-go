package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/oudentabetai/pterodactyl-go/pterodactyl"
	"github.com/oudentabetai/pterodactyl-go/storage"
	"github.com/oudentabetai/pterodactyl-go/utils"
)

func HelpCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This is the Help Message.",
		},
	})
	if err != nil {
		return
	}
}

func ServersCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 1. 保留応答
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Println("InteractionRespond error:", err)
		return
	}
	var serverIDs []string
	for _, role := range i.Member.Roles {
		serverIDs = append(serverIDs, storage.ConfigMgr.GetServerID(role)...)
	}

	// 2. サーバー情報取得
	response := pterodactyl.GetServers(s)
	if response == nil {
		log.Println("Error fetching servers:", err)

		// ユーザーにエラーを通知して終了（「考えています」状態を解除）
		errMsg := "❌ サーバー情報の取得中にエラーが発生しました。"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}
	defer response.Body.Close()
	json := utils.GetServerJson(response)

	serverHasAccess := make([]utils.Server, 0)
	for _, server := range json.Data {
		for _, id := range serverIDs {
			attributes := string(rune(server.Attributes.ID))
			if attributes == id {
				log.Printf("User has access to server: %s (ID: %s)", server.Attributes.Name, server.Attributes.ID)
				serverHasAccess = append(serverHasAccess, server)
			}

		}
	}

	if serverHasAccess == nil {
		log.Println("Error creating embed from server response")
		errMsg := "❌ サーバー情報の整形中にエラーが発生しました。"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}
	// 3. 正常終了（Embedを表示）
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{responseEmbed},
	})

	if err != nil {
		log.Println("Embed送信エラー:", err)
	}
}

func ServerCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Println("InteractionRespond error:", err)
		return
	}

	var action, identifier, id string
	for _, option := range i.ApplicationCommandData().Options {
		if option.Name == "action" {
			action = option.StringValue()
		}
		if option.Name == "server_id" {
			identifier = option.StringValue()
		}
		if option.Name == "id" {
			id = option.StringValue()
		}
	}
	if (id == "" || identifier == "") && action == "" {
		errMsg := "❌ コマンドの形式が正しくありません。"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

}
