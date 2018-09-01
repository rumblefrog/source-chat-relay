package bot

import (
	"log"

	"../helper"
	"github.com/bwmarrin/discordgo"
)

// Session - Initiated session
var Session *discordgo.Session

// InitBot - Starts the bot routine
func InitBot() {
	Session, err := discordgo.New("Bot" + helper.Conf.Token)

	if err != nil {
		log.Panic("Unable to initiate bot session")
	}

	Session.AddHandler(messageCreate)

	err = Session.Open()

	if err != nil {
		log.Panic("Unable to open bot connection")
	}

	log.Println("Bot is now running")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Prevent loops
	if m.Author.ID == s.State.User.ID {
		return
	}

}
