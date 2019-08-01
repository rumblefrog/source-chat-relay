package bot

import (
	"database/sql"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/source-chat-relay/server/entity"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/sirupsen/logrus"
)

type ChannelCmdType uint8

const (
	Receive ChannelCmdType = iota
	Send
)

func channelCommandRoute(ctx *exrouter.Context) {
	cType := ctx.Args.Get(0)

	if cType == "receivechannel" {
		channelCommand(ctx, Receive)
	} else if cType == "sendchannel" {
		channelCommand(ctx, Send)
	}
}

func channelCommand(ctx *exrouter.Context, cmdType ChannelCmdType) {
	var (
		dChannel *discordgo.Channel
		err      error
	)

	id := ctx.Args.Get(1)
	channel := strings.TrimSpace(ctx.Args.After(2))

	pid, ok := ParseChannel(id)

	if ok {
		id = pid
		dChannel, err = ctx.Ses.Channel(id)

		if err != nil {
			logrus.WithField("error", err).Warn("Unable to fetch channel")

			ctx.Reply("Unable to fetch channel")

			return
		}
	}

	tEntity, err := entity.GetEntity(id)

	if err == sql.ErrNoRows && channel != "" {
		tEntity = &entity.Entity{
			ID:        id,
			CreatedAt: time.Now(),
		}

		if cmdType == Receive {
			tEntity.ReceiveChannels = entity.ParseDelimitedChannels(channel)
		} else if cmdType == Send {
			tEntity.SendChannels = entity.ParseDelimitedChannels(channel)
		}

		if ok {
			tEntity.DisplayName = dChannel.Name
		}

		err = tEntity.Insert()

		if err != nil {
			ctx.Reply("Unable to create entity")
			return
		}
	} else if err != nil {
		logrus.WithField("error", err).Warn("Unable to fetch entity")

		ctx.Reply("Unable to fetch entity")

		return
	} else if channel != "" {
		if cmdType == Receive {
			err = tEntity.SetReceiveChannels(entity.ParseDelimitedChannels(channel))
		} else if cmdType == Send {
			err = tEntity.SetSendChannels(entity.ParseDelimitedChannels(channel))
		}

		if ok {
			tEntity.DisplayName = dChannel.Name
		}

		if err != nil {
			logrus.WithField("error", err).Warn("Unable to update entity")

			ctx.Reply("Unable to update entity")

			return
		}
	}

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, tEntity.Embed())
}
