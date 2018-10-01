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
		if cmdType == Receive {
			entity = &repoEntity.Entity{
				ID:              id,
				Type:            eType,
				ReceiveChannels: repoEntity.ParseChannels(channel),
				CreatedAt:       time.Now(),
			}
		} else if cmdType == Send {
			entity = &repoEntity.Entity{
				ID:           id,
				Type:         eType,
				SendChannels: repoEntity.ParseChannels(channel),
				CreatedAt:    time.Now(),
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
		err = entity.SetReceiveChannels(repoEntity.ParseChannels(channel))

		if err != nil {
			log.WithField("error", err).Warn("Unable to update entity")

			ctx.Reply("Unable to update entity")

			return
		}
	}

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, entity.Embed())
}
