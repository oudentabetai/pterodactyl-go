package pterodactylbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/oudentabetai/pterodactyl-go/discord"
)

func Register(s *discordgo.Session) {
	s.AddHandler(discord.OnMessageCreate)
	s.AddHandler(discord.OnInteractionCreate)
}
