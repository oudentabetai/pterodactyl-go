package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/oudentabetai/pterodactyl-go/pterodactyl"
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

	// 2. サーバー情報取得
	response, err := pterodactyl.GetServers(s, *i.Member)
	if err != nil {
		log.Println("Error fetching servers:", err)

		// ユーザーにエラーを通知して終了（「考えています」状態を解除）
		errMsg := "❌ サーバー情報の取得中にエラーが発生しました。"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	// 3. 正常終了（Embedを表示）
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{response},
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
		if option.Name == "ServerId" {
			identifier = option.StringValue()
		}
		if option.Name == "Id" {
			id = option.StringValue()
		}
	}
	if (id == "" || identifier == "") && action == "" {
		errMsg := "❌ コマンドの形式が正しくありません。例: /server <action> <serverIdentifier>"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

}
