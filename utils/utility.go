package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/oudentabetai/pterodactyl-go/pterodactyl"
	"github.com/oudentabetai/pterodactyl-go/storage"
)

type PteroResponse[T any] struct {
	Data []T `json:"data"`
}

type Server struct {
	Attributes struct {
		ID         int    `json:"id"`
		UUID       string `json:"uuid"`
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
		Status     string `json:"status"`
	} `json:"attributes"`
}

type User struct {
	Attributes struct {
		ID     int    `json:"id"`
		Name   string `json:"username"`
		Email  string `json:"email"`
		IsRoot bool   `json:"is_root"`
	} `json:"attributes"`
}

type Status string

const (
	StatusRunning  Status = "running"
	StatusOffline  Status = "offline"
	StatusStarting Status = "starting"
	StatusStopping Status = "stopping"
)

func (s Status) ToJapanese() string {
	switch s {
	case StatusRunning:
		return "起動中"
	case StatusOffline:
		return "停止"
	case StatusStarting:
		return "起動処理中"
	case StatusStopping:
		return "停止処理中"
	default:
		return string(s)
	}
}

func GetRoles(roleIDs []string) []string {
	var serverIDs []string
	for _, roleID := range roleIDs {
		serverIDs = append(serverIDs, storage.ConfigMgr.GetServerID(roleID)...)
	}
	return serverIDs
}

func GetServerJson(resp *http.Response) *PteroResponse[Server] {
	var serverres PteroResponse[Server]
	if err := json.NewDecoder(resp.Body).Decode(&serverres); err != nil {
		log.Printf("Decode error: %v", err)
		return nil
	}
	return &serverres
}

func ListServers(resp *http.Response) *discordgo.MessageEmbed {
	serverres := GetServerJson(resp)
	if serverres == nil {
		return &discordgo.MessageEmbed{
			Title:       "サーバーリスト",
			Description: "サーバー情報の取得に失敗しました。",
			Color:       0xff0000,
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "サーバーリスト",
		Description: "取得したサーバーのリストです。",
		Color:       0x00ff00,
	}

	for i, sv := range serverres.Data {
		status := Status(pterodactyl.GetServerStatus(sv.Attributes.Identifier))

		statusEmoji := "⚪"
		switch status {
		case StatusRunning:
			statusEmoji = "🟢"
		case StatusOffline:
			statusEmoji = "🔴"
		case StatusStarting:
			statusEmoji = "🟡"
		}

		fieldValue := fmt.Sprintf("🆔 **ID**: %d\n🔑 **Identifier**: `%s`\n📊 **Status**: %s %s",
			i+1,
			sv.Attributes.Identifier,
			statusEmoji,
			status.ToJapanese(),
		)

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🏠 " + sv.Attributes.Name,
			Value:  fieldValue,
			Inline: false,
		})
	}
	return embed
}

func ListUsers(resp *http.Response) string {
	var userres PteroResponse[User]
	if err := json.NewDecoder(resp.Body).Decode(&userres); err != nil {
		log.Printf("Decode error: %v", err)
		return "ユーザーの情報を取得できませんでした"
	}
	result := "ユーザー情報:\n"
	for _, user := range userres.Data {
		result += "ID: " + string(rune(user.Attributes.ID)) + "\n"
	}
	return result
}
