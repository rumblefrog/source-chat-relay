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
		log.WithField("error", err).Fatal("Unable to initiate bot session")
	}

	session.AddHandler(ready)

	err = session.Open()

	if err != nil {
		log.WithField("error", err).Fatal("Unable to open bot session")
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
				Hostname:   "Discord",
				ClientName: m.Author.Username,
				ClientID:   m.Author.ID,
			}

			protocol.NetManager.Router <- message
		}
	})

	RelayBot = &DiscordBot{
		Session: session,
	}

	router.Group(func(r *exrouter.Route) {
		r.Cat("configuration")

		r.Use(Auth)

		r.On("receivechannel", ChannelCommandRoute).Desc("Get/Set the receive relay channel of this ID/TextChannel").Alias("rc")

		r.On("sendchannel", ChannelCommandRoute).Desc("Get/Set the send relay channel of this ID/TextChannel").Alias("sc")
	})

	router.Group(func(r *exrouter.Route) {
		r.Cat("other")

		r.On("ping", func(ctx *exrouter.Context) {
			ctx.Reply("pong")
		}).Desc("Responds with pong").Cat("other")
	})
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	go RelayBot.Listen()

	log.WithFields(log.Fields{
		"Username":    event.User.Username,
		"Guild Count": len(event.Guilds),
	}).Info("Bot is now running")
}
