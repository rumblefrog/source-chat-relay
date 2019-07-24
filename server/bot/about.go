package bot

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/relay"
)

func aboutCommand(ctx *exrouter.Context) {
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
			&discordgo.MessageEmbedField{
				Name:  "Incoming Traffic",
				Value: relay.Instance.Statistics.Incoming.String(),
			},
			&discordgo.MessageEmbedField{
				Name:  "Outgoing Traffic",
				Value: relay.Instance.Statistics.Outgoing.String(),
			},
		},
	})
}
