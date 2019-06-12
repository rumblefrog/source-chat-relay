package bot

import (
	"github.com/Necroforger/dgrouter"
	"github.com/Necroforger/dgrouter/exrouter"
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/protocol"

	repoEntity "github.com/rumblefrog/source-chat-relay/server/repositories/entity"
)

type DiscordBot struct {
	Session *discordgo.Session
}

var RelayBot *DiscordBot

func init() {
	session, err := discordgo.New("Bot " + config.Conf.Bot.Token)

	if err != nil {
		log.WithField("error", err).Fatal("Unable to initiate bot session")
	}

	RelayBot = &DiscordBot{
		Session: session,
	}

	session.AddHandler(ready)

	err = session.Open()

	if err != nil {
		log.WithField("error", err).Fatal("Unable to open bot session")
	}

	router := exrouter.New()

	session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot && !config.Conf.Bot.ListenToBots {
			return
		}

		err := router.FindAndExecute(session, "r/", session.State.User.ID, m.Message)

		if err == dgrouter.ErrCouldNotFindRoute {

			relayChannel, err := repoEntity.GetEntity(m.ChannelID, repoEntity.Channel)

			if err != nil {
				return
			}

			channel, err := session.Channel(m.ChannelID)

			if err != nil {
				return
			}

			message := &protocol.Message{
				Overwrite: &protocol.OverwriteData{
					SendChannels: relayChannel.SendChannels,
				},
				Hostname:   CapitalChannelName(channel),
				ClientName: m.Author.Username,
				ClientID:   m.Author.ID,
				Content:    TransformMentions(session, m.ChannelID, m.Content),
			}

			protocol.NetManager.Router <- message
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
	go RelayBot.Listen()

	log.WithFields(log.Fields{
		"Username":    event.User.Username,
		"Guild Count": len(event.Guilds),
	}).Info("Bot is now running")
}
