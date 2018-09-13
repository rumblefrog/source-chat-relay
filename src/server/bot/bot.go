package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/src/server/helper"
)

var Session *discordgo.Session

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
