package discord

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/oudentabetai/pterodactyl-go/pterodactyl"
	"github.com/oudentabetai/pterodactyl-go/storage"
	"github.com/oudentabetai/pterodactyl-go/utils"
)

var (
	suffix   string = "!!"
	OWNER_ID string = "967088187405107220"
)

type SessionManager interface {
	InitializeSession(token string) *discordgo.Session
}

type DiscordSessionManager struct{}

func safeSendEmbed(s *discordgo.Session, channelID string, embed *discordgo.MessageEmbed) {
	if embed == nil {
		embed = &discordgo.MessageEmbed{
			Title:       "サーバーリスト",
			Description: "サーバー情報の整形に失敗しました。",
			Color:       0xff0000,
		}
	}

	if _, err := s.ChannelMessageSendEmbed(channelID, embed); err != nil {
		log.Printf("failed to send embed: %v", err)
	}
}

func (d *DiscordSessionManager) InitializeSession(token string) *discordgo.Session {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Discordセッションの作成に失敗: %v", err)
	}
	return dg
}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered in OnMessageCreate: %v", r)
		}
	}()

	if m == nil || m.Author == nil {
		return
	}

	u := m.Author
	if !u.Bot {
		if strings.HasPrefix(m.Content, suffix) {
			log.Printf("Received command: %s from user: %s", m.Content, u.Username)
		}
		if strings.HasPrefix(m.Content, suffix+"user") {
			s.ChannelMessageSend(m.ChannelID, pterodactyl.GetUser())
		}
		if strings.HasPrefix(m.Content, suffix+"servers") {
			resp := pterodactyl.GetServers(s)
			if resp == nil {
				s.ChannelMessageSend(m.ChannelID, "サーバー情報の取得に失敗しました。")
				return
			}
			defer resp.Body.Close()
			embed := utils.ListServers(resp)
			safeSendEmbed(s, m.ChannelID, embed)
			return
		}
		if strings.HasPrefix(m.Content, suffix+"setrole") {
			if m.Author.ID != OWNER_ID {
				log.Print("User does not have permission to set role: " + OWNER_ID + " vs " + m.Author.ID)
				s.ChannelMessageSend(m.ChannelID, "このコマンドを使用する権限がありません。")
				return
			}
			fields := strings.Fields(m.Content)
			if len(fields) == 3 {
				s.ChannelMessageSend(m.ChannelID, storage.ConfigMgr.SetRole(fields[1], fields[2]))
			} else {
				s.ChannelMessageSend(m.ChannelID, "コマンドの形式が正しくありません。例: !!setrole <roleID> <serverID>")
			}
		}
		if strings.HasPrefix(m.Content, suffix+"server") {
			fields := strings.Fields(m.Content)
			if len(fields) == 3 {
				s.ChannelMessageSend(m.ChannelID, pterodactyl.PowerServer(fields[1], fields[2]))
			} else {
				s.ChannelMessageSend(m.ChannelID, "コマンドの形式が正しくありません。例: !!server <action> <serverIdentifier>")
			}
		}
	}
}
