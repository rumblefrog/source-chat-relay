package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/src/server/helper"
	"github.com/rumblefrog/source-chat-relay/src/server/protocol"
)

type DiscordBot struct {
	Session *discordgo.Session
	Data    chan *protocol.Message
}

var RelayBot *DiscordBot

func InitBot() {
	session, err := discordgo.New("Bot" + helper.Conf.Bot.Token)

	if err != nil {
		log.Panic("Unable to initiate bot session")
	}

	session.AddHandler(messageCreate)

	err = session.Open()

	if err != nil {
		log.Panic("Unable to open bot connection")
	}

	RelayBot = &DiscordBot{
		Session: session,
		Data:    make(chan *protocol.Message),
	}

	log.Println("Bot is now running")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Prevent loops
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Send to router directly aftering constructing struct
}
