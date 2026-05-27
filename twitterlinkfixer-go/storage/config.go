package storage

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DISCORD_BOT_TOKEN string
	APPLICATION_ID    string
	TEST_GUILD_ID     string
	OWNER_ID          string
	LOG_CHANNEL_ID    string
}

var Envs = loadEnvs()

func loadEnvs() Env {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return Env{
		APPLICATION_ID:    os.Getenv("APPLICATION_ID"),
		DISCORD_BOT_TOKEN: os.Getenv("DISCORD_BOT_TOKEN"),
		TEST_GUILD_ID:     os.Getenv("TEST_GUILD_ID"),
		OWNER_ID:          os.Getenv("OWNER_ID"),
		LOG_CHANNEL_ID:    os.Getenv("LOG_CHANNEL_ID"),
	}
}
