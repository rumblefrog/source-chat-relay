package bot

import (
	"github.com/Necroforger/dgrouter"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/protocol"
	"github.com/rumblefrog/source-chat-relay/server/relay"
)

var RelayBot *discordgo.Session

func Initialize() {
	session, err := discordgo.New("Bot " + config.Config.Bot.Token)

	if err != nil {
		logrus.WithField("error", err).Fatal("Unable to initiate bot session")
	}

	RelayBot = session

	session.AddHandler(ready)

	session.State.TrackEmojis = false
	session.State.TrackPresences = false
	session.State.TrackVoice = false

	err = session.Open()

	if err != nil {
		logrus.WithField("error", err).Fatal("Unable to open bot session")
	}

	router := exrouter.New()

	session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot && !config.Config.Bot.ListenToBots {
			return
		}

		err := router.FindAndExecute(session, "r/", session.State.User.ID, m.Message)

		if err == dgrouter.ErrCouldNotFindRoute {
			channel, err := session.Channel(m.ChannelID)

			if err != nil {
				return
			}

			transformed, err := m.ContentWithMoreMentionsReplaced(session)

			if err != nil {
				transformed = m.Content
			}

			message := &protocol.ChatMessage{
				BaseMessage: protocol.BaseMessage{
					Type:       protocol.MessageChat,
					SenderID:   m.ChannelID,
					EntityName: CapitalChannelName(channel),
				},
				IDType:   protocol.IdentificationDiscord,
				ID:       m.Author.ID,
				Username: m.Author.Username,
				Message:  transformed,
			}

			relay.Instance.Router <- message
		}
	})

	router.Group(func(r *exrouter.Route) {
		r.Cat("configuration")

		r.Use(Auth)

		r.On("receivechannel", ChannelCommandRoute).Desc("Get/Set the receive relay channel of this ID/TextChannel").Alias("rc")

		r.On("sendchannel", ChannelCommandRoute).Desc("Get/Set the send relay channel of this ID/TextChannel").Alias("sc")

		r.On("entities", EntitiesCMD).Desc("Fetch all entities (of certain type)").Alias("e")
	})

	router.Group(func(r *exrouter.Route) {
		r.Cat("other")

		r.On("about", func(ctx *exrouter.Context) {
			ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{
					Name:    "Fishy!",
					URL:     "https://github.com/rumblefrog",
					IconURL: "https://avatars2.githubusercontent.com/u/6960234?s=32",
				},
				Color: 0x3395D6,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:  "SCR Version",
						Value: config.SCRVER,
					},
					&discordgo.MessageEmbedField{
						Name:  "Repository",
						Value: "https://github.com/rumblefrog/source-chat-relay/",
					},
				},
			})
		}).Desc("Version and source information").Alias("info")

		r.On("ping", func(ctx *exrouter.Context) {
			ctx.Reply("pong")
		}).Desc("Responds with pong")
	})
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	go Listen()

	logrus.WithFields(logrus.Fields{
		"Username":    event.User.Username,
		"Guild Count": len(event.Guilds),
	}).Info("Bot is now running")
}
