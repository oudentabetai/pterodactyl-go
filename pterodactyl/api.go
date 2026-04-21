package pterodactyl

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/oudentabetai/pterodactyl-go/storage"
)

var (
	APPLICATION_API_KEY string
	CLIENT_API_KEY      string
	BASE_URL            = "https://web.ofton.dev/api/"
)

func SetAPIKeys(applicationAPIKey, clientAPIKey string) {
	APPLICATION_API_KEY = applicationAPIKey
	CLIENT_API_KEY = clientAPIKey
}

type PteroResponse[T any] struct {
	Data []T `json:"data"`
}

type User struct {
	Attributes struct {
		ID     int    `json:"id"`
		Name   string `json:"username"`
		Email  string `json:"email"`
		IsRoot bool   `json:"is_root"`
	} `json:"attributes"`
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

type ResourceResponse struct {
	Attributes struct {
		CurrentState string `json:"current_state"`
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

func Fetch(URL string, apiKey string) *http.Response {
	log.Print("Fetching :" + URL)
	req, err := http.NewRequest(http.MethodGet, URL, nil)

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/vnd.pterodactyl.v1+json")

	if err != nil {
		log.Printf("Request error: %v", err)
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do error: %v", err)
	}

	return resp
}

func GetUser() string {
	if APPLICATION_API_KEY == "" {
		log.Print("APPLICATION_API_KEY is empty")
		return "APIキーが未設定です"
	}

	resp := Fetch(BASE_URL+"application/users", APPLICATION_API_KEY)
	if resp == nil {
		return "API呼び出しに失敗しました"
	}
	defer resp.Body.Close()

	var res PteroResponse[User]
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Printf("Decode error: %v", err)
		return "Decode failed"
	}

	var result string
	for _, user := range res.Data {
		result += user.Attributes.Name + ", "
	}
	return result
}

func GetStatus(s *discordgo.Session, m *discordgo.MessageCreate) {
	if APPLICATION_API_KEY == "" || CLIENT_API_KEY == "" {
		log.Print("API keys are empty")
		return
	}

	resp := Fetch(BASE_URL+"application/servers", APPLICATION_API_KEY)
	if resp == nil {
		log.Print("Failed to fetch server list")
		return
	}
	defer resp.Body.Close()

	var serverres PteroResponse[Server]
	if err := json.NewDecoder(resp.Body).Decode(&serverres); err != nil {
		log.Printf("Decode error: %v", err)
	}
	servers := serverres.Data
	for i := range servers {
		statusResp := Fetch(BASE_URL+"client/servers/"+servers[i].Attributes.Identifier+"/resources", CLIENT_API_KEY)
		if statusResp == nil {
			continue
		}
		var resourceres ResourceResponse
		if err := json.NewDecoder(statusResp.Body).Decode(&resourceres); err == nil {
			servers[i].Attributes.Status = resourceres.Attributes.CurrentState
		}
		statusResp.Body.Close()

		servers[i].Attributes.Status = resourceres.Attributes.CurrentState
	}

	if m.Member == nil {
		log.Print("Member is nil; cannot resolve roles")
		return
	}

	roles := m.Member.Roles
	log.Print("Member roles: " + strings.Join(roles, ", "))

	matchedServers := make([]Server, 0)
	added := make(map[string]bool)
	for _, roleID := range roles {
		log.Print("Checking role ID: " + roleID)
		serverID := storage.ConfigMgr.GetServerID(roleID)
		if serverID == "" {
			continue
		}

		for i := range servers {
			log.Print("Checking server: " + servers[i].Attributes.Identifier + " and " + serverID)
			if servers[i].Attributes.Identifier == serverID && !added[serverID] {
				matchedServers = append(matchedServers, servers[i])
				added[serverID] = true
			}
		}
	}

	if len(matchedServers) > 0 {
		sendServerListEmbed(s, m, matchedServers)
	} else {
		s.ChannelMessageSend(m.ChannelID, "一致するサーバーが見つかりませんでした。")
	}
}

func ServerManager(m *discordgo.MessageCreate, action, serverIdentifier string) string {
	for _, roleID := range m.Member.Roles {
		if storage.ConfigMgr.GetServerID(roleID) == serverIdentifier {
			log.Print("User has permission to " + action + " server: " + serverIdentifier)
			var signal string
			switch action {
			case "start":
				signal = "start"
			case "stop":
				signal = "stop"
			case "restart":
				signal = "restart"
			default:
				return "不明なアクションです。使用可能なアクション: start, stop, restart"
			}
			req, err := http.NewRequest(http.MethodPost, BASE_URL+"client/servers/"+serverIdentifier+"/power", strings.NewReader(`{"signal":"`+signal+`"}`))
			req.Header.Set("Authorization", "Bearer "+CLIENT_API_KEY)
			req.Header.Set("Accept", "application/vnd.pterodactyl.v1+json")
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				log.Printf("Request creation error: %v", err)
				return "サーバー " + action + " リクエストの作成に失敗しました"
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Do error: %v", err)
				return "サーバーの " + action + " リクエストの送信に失敗しました"
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNoContent {
				return "サーバーの " + action + " を開始しました: " + serverIdentifier
			} else {
				log.Printf("Unexpected status code: %d", resp.StatusCode)
				return fmt.Sprintf("サーバーの "+action+" に失敗しました (Status Code: %d)", resp.StatusCode)
			}
		}
	}

	return "このサーバーを " + action + " する権限がありません"
}

func sendServerListEmbed(s *discordgo.Session, m *discordgo.MessageCreate, servers []Server) {
	embed := &discordgo.MessageEmbed{
		Title:       "🎮 サーバー稼働状況一覧",
		Description: "現在管理中のサーバーリストです。",
		Color:       0x00ff00,
		Fields:      []*discordgo.MessageEmbedField{},
	}

	for i, sv := range servers {
		attr := sv.Attributes
		status := Status(attr.Status)

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
			attr.Identifier,
			statusEmoji,
			status.ToJapanese(),
		)

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "🏠 " + attr.Name,
			Value:  fieldValue,
			Inline: false,
		})
	}

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
