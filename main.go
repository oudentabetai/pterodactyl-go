package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/oudentabetai/pterodactyl-go/discord"
	"github.com/oudentabetai/pterodactyl-go/pterodactyl"
	"github.com/oudentabetai/pterodactyl-go/storage"
	"github.com/oudentabetai/pterodactyl-go/utils"
)

var (
	GuildID string
	dgs     *discordgo.Session
)

func Env_load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	sessionManager := &discord.DiscordSessionManager{}
	Env_load()
	utils.OWNER_ID = os.Getenv("OWNER_ID")
	pterodactyl.SetAPIKeys(
		os.Getenv("PANEL_API_TOKEN"),
		os.Getenv("PANEL_CLIENT_TOKEN"),
	)
	discord.GetLogChannelID(os.Getenv("LOG_CHANNEL_ID"))
	if err := storage.ConfigMgr.Load(); err != nil {
		log.Fatalf("設定ファイルの読み込みに失敗: %v", err)
	}
	dgs = sessionManager.InitializeSession(os.Getenv("DISCORD_BOT_TOKEN"))
	dgs.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsAll | discordgo.PermissionSendMessages
	if err := dgs.Open(); err != nil {
		log.Fatalf("Discordセッションのオープンに失敗: %v", err)
	}
	dgs.AddHandler(discord.OnMessageCreate)
	dgs.AddHandler(discord.OnInteractionCreate)
	defer dgs.Close()
	log.Println("ボットが起動しました。Ctrl+Cで終了します。")

	//deleteAllGlobalCommands(dgs, os.Getenv("APPLICATION_ID"))
	discord.SyncCommands(dgs, "", os.Getenv("APPLICATION_ID"))
	waitForExitSignal()
}

func waitForExitSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func GetOwnerID() string {
	return os.Getenv("OWNER_ID")
}

func deleteAllGlobalCommands(s *discordgo.Session, appID string) {
	_, err := s.ApplicationCommandBulkOverwrite(appID, "", []*discordgo.ApplicationCommand{})

	if err != nil {
		log.Printf("グローバルコマンドの削除に失敗しました: %v", err)
		return
	}
	log.Println("すべてのグローバルコマンドを削除しました。")
}
