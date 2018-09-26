package bot

import (
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/src/server/helper"
	"github.com/rumblefrog/source-chat-relay/src/server/protocol"
)

type DiscordBot struct {
	Session       *discordgo.Session
	Data          chan *protocol.Message
	RelayChannels []*RelayChannel
}

type RelayChannel struct {
	Channel         *discordgo.Channel
	ReceiveChannels []int
	SendChannels    []int
}

var RelayBot *DiscordBot

// TODO: rename to init
func InitBot() {
	session, err := discordgo.New("Bot" + helper.Conf.Bot.Token)

	if err != nil {
		log.Fatal("Unable to initiate bot session")
	}

	session.AddHandler(messageCreate)

	err = session.Open()

	if err != nil {
		log.Fatal("Unable to open bot connection")
	}

	RelayBot = &DiscordBot{
		Session: session,
		Data:    make(chan *protocol.Message),
	}

	log.Info("Bot is now running")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Prevent loops
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Create fake Client for Message

	// Send to router directly aftering constructing struct
}
