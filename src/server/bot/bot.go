package bot

import (
	"log"

	"../helper"
	"github.com/bwmarrin/discordgo"
)

// Bot - Initiated session
var Bot *discordgo.Session

// InitBot - Starts the bot routine
func InitBot() {
	Bot, err := discordgo.New("Bot" + helper.Conf.Token)

	if err != nil {
		log.Panic("Unable to initiate bot session")
	}

	Bot.AddHandler(messageCreate)

	err = Bot.Open()

	if err != nil {
		log.Panic("Unable to open bot connection")
	}

	log.Println("Bot is now running")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

}
