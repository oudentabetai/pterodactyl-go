package discord

import (
  "github.com/bwmarrin/discordgo"
)

func HelpCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
  i.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
    Type: discordgo.InteractionResponseChannelMessageWithSource,
    Content:  &discordgo.InteractionResponseData{
					Content: "hi",
      }
  })
}

func ServersCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
  i.InteractionRespond(i.Interaction, &discordgo.interactionResponse{
    Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
    })
  res, err := pterodactyl.GetServers()
  if err {
    log.Panic(err)
    }
  
}
