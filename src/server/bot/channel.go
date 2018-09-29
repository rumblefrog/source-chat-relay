package bot

import (
	"database/sql"
	"time"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/rumblefrog/source-chat-relay/src/server/database"
	log "github.com/sirupsen/logrus"
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
	cType := database.Server

	pid, ok := ParseChannel(id)

	if ok {
		cType = database.Channel
		id = pid
	}

	log.WithFields(log.Fields{
		"ID":      id,
		"Channel": channel,
	}).Debug()

	entity, err := database.FetchEntity(id)

	if err == sql.ErrNoRows && channel != "" {
		if cmdType == Receive {
			entity = &database.Entity{
				ID:              id,
				Type:            cType,
				ReceiveChannels: database.ParseChannels(channel),
				SendChannels:    []int{},
				CreatedAt:       time.Now(),
			}
		} else if cmdType == Send {
			entity = &database.Entity{
				ID:              id,
				Type:            cType,
				ReceiveChannels: []int{},
				SendChannels:    database.ParseChannels(channel),
				CreatedAt:       time.Now(),
			}
		}

		_, err = entity.CreateEntity()

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
			entity.ReceiveChannels = database.ParseChannels(channel)
		} else if cmdType == Send {
			entity.SendChannels = database.ParseChannels(channel)
		}

		_, err = entity.UpdateChannels()

		if err != nil {
			log.WithField("error", err).Warn("Unable to update entity")

			ctx.Reply("Unable to update entity")

			return
		}
	}

	DisplayEntity(ctx, entity, "Entity Descriptor")
}
