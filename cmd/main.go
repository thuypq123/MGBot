package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"discord_bot/internal/usecase"
	"discord_bot/internal/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {
	envDiscordToken, err := utils.GetSecretAws()
	if err != nil {
		fmt.Println("Error getting AWS secret: ", err)
		return
	}
	if envDiscordToken.BotToken == "" {
		fmt.Println("DISCORD_BOT_TOKEN is not set")
		return
	}
	fmt.Println("DISCORD_BOT_TOKEN: ", utils.MaskSensitiveString(envDiscordToken.BotToken))
	discord, err := createDiscordSession(envDiscordToken.BotToken)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	err = openDiscordConnection(discord)
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}
	discord.AddHandler(usecase.HandleMessageVoice)

	waitForExitSignal()

	discord.Close()
}

func createDiscordSession(token string) (*discordgo.Session, error) {
	return discordgo.New("Bot " + token)
}

func openDiscordConnection(discord *discordgo.Session) error {
	return discord.Open()
}

func waitForExitSignal() {
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
