package bot

import (
	"github.com/Necroforger/dgrouter"
	"github.com/Necroforger/dgrouter/exrouter"
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/src/server/helper"
	"github.com/rumblefrog/source-chat-relay/src/server/protocol"
)

type DiscordBot struct {
	Session       *discordgo.Session
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

	err = session.Open()

	if err != nil {
		log.Fatal("Unable to open bot connection", err)
	}

	router := exrouter.New()

	session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		err := router.FindAndExecute(session, "r/", session.State.User.ID, m.Message)

		if err == dgrouter.ErrCouldNotFindRoute {

			relayChannel := RelayBot.GetRelayChannel(m.ChannelID)

			if relayChannel == nil {
				return
			}

			message := &protocol.Message{
				Overwrite: &protocol.OverwriteData{
					SendChannels: relayChannel.SendChannels,
				},
				ClientName: m.Author.Username,
				ClientID:   m.Author.ID,
			}

			protocol.NetManager.Router <- message
		}
	})
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	RelayBot = &DiscordBot{
		Session: s,
	}

	log.WithFields(log.Fields{
		"Username":    event.User.Username,
		"Session ID":  event.SessionID,
		"Guild Count": len(event.Guilds),
	}).Info("Bot is now running")
}
