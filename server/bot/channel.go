package bot

import (
	"database/sql"
	"time"

	"github.com/rumblefrog/source-chat-relay/server/repositories/entity"

	"github.com/Necroforger/dgrouter/exrouter"
	log "github.com/sirupsen/logrus"

	repoEntity "github.com/rumblefrog/source-chat-relay/server/repositories/entity"
)

type ChannelCmdType int

const (
	Receive ChannelCmdType = iota
	Send
)

func ChannelCommandRoute(ctx *exrouter.Context) {
	cType := ctx.Args.Get(0)

	if cType == "receivechannel" {
		ChannelCommand(ctx, Receive)
	} else if cType == "sendchannel" {
		ChannelCommand(ctx, Send)
	}
}

func ChannelCommand(ctx *exrouter.Context, cmdType ChannelCmdType) {
	id := ctx.Args.Get(1)
	channel := ctx.Args.After(2)
	eType := repoEntity.Server

	pid, ok := ParseChannel(id)

	if ok {
		eType = repoEntity.Channel
		id = pid
	}

	entity, err := entity.GetEntity(id, eType)

	if err == sql.ErrNoRows && channel != "" {
		entity = &repoEntity.Entity{
			ID:        id,
			Type:      eType,
			CreatedAt: time.Now(),
		}

		if cmdType == Receive {
			entity.ReceiveChannels = repoEntity.ParseChannels(channel)
		} else if cmdType == Send {
			entity.SendChannels = repoEntity.ParseChannels(channel)
		}

		if eType == repoEntity.Channel {
			dChannel, err := ctx.Ses.Channel(id)

			if err != nil {
				log.WithField("error", err).Warn("Unable to fetch channel for insertion")
			} else {
				entity.DisplayName = dChannel.Name
			}
		}

		err = entity.Insert()

		if err != nil {
			ctx.Reply("Unable to create entity")
			return
		}
	} else if err != nil {
		log.WithField("error", err).Warn("Unable to fetch entity")

		ctx.Reply("Unable to fetch entity")

		return
	} else {
		if cmdType == Receive {
			err = entity.SetReceiveChannels(repoEntity.ParseChannels(channel))
		} else if cmdType == Send {
			err = entity.SetSendChannels(repoEntity.ParseChannels(channel))
		}
		if err != nil {
			log.WithField("error", err).Warn("Unable to update entity")

			ctx.Reply("Unable to update entity")

			return
		}
	}

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, entity.Embed())
}
