package bot

import (
	"strings"

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

	session.AddHandler(func(session *discordgo.Session, m *discordgo.MessageCreate) {
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

			displayname := m.Author.Username

			member, err := session.GuildMember(m.GuildID, m.Author.ID)

			if err == nil && len(member.Nick) != 0 {
				displayname = member.Nick
			}

			message := &protocol.ChatMessage{
				BaseMessage: protocol.BaseMessage{
					Type:       protocol.MessageChat,
					SenderID:   m.ChannelID,
					EntityName: strings.Title(channel.Name),
				},
				IDType:   protocol.IdentificationDiscord,
				ID:       m.Author.ID,
				Username: displayname,
				Message:  transformed,
			}

			relay.Instance.Router <- message
		}
	})

	router.Group(func(r *exrouter.Route) {
		r.Cat("configuration")

		r.Use(Auth)

		r.On("receivechannel", channelCommandRoute).Desc("Get/Set the receive relay channel of this ID/TextChannel").Alias("rc")

		r.On("sendchannel", channelCommandRoute).Desc("Get/Set the send relay channel of this ID/TextChannel").Alias("sc")

		r.On("deleteentity", deleteCommand).Desc("Delete an entity from database").Alias("delete", "del")

		r.On("entities", entitiesCommand).Desc("Fetch all entities (of certain type)").Alias("e")
	})

	router.Group(func(r *exrouter.Route) {
		r.Cat("message")

		r.Use(Auth)

		r.On("event", eventCommand).Desc("Send an event message")
	})

	router.Group(func(r *exrouter.Route) {
		r.Cat("other")

		r.On("about", aboutCommand).Desc("Version, source, stats information").Alias("info")

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
