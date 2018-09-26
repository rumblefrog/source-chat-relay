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
	ChannelID       string
	ReceiveChannels []int
	SendChannels    []int
}

var RelayBot *DiscordBot

func init() {
	session, err := discordgo.New("Bot " + helper.Conf.Bot.Token)

	if err != nil {
		log.Fatal("Unable to initiate bot session")
	}

	session.AddHandler(ready)
	session.AddHandler(messageCreate)

	err = session.Open()

	if err != nil {
		log.Fatal("Unable to open bot connection", err)
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	RelayBot = &DiscordBot{
		Session: s,
		Data:    make(chan *protocol.Message),
	}

	log.WithFields(log.Fields{
		"Username":    event.User.Username,
		"Session ID":  event.SessionID,
		"Guild Count": len(event.Guilds),
	}).Info("Bot is now running")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Prevent loops
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Create fake Client for Message

	// Send to router directly aftering constructing struct
}
