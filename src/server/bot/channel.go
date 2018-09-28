package bot

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/rumblefrog/source-chat-relay/src/server/database"
	log "github.com/sirupsen/logrus"
)

func ChannelCommand(ctx *exrouter.Context) {
	id := ctx.Args.Get(1)

	pid, ok := ParseChannel(id)

	if ok {
		id = pid
	}

	log.Debug(id)

	entity, err := database.FetchEntity(id)

	if err != nil {
		log.WithField("error", err).Warn("Unable to fetch entity")

		ctx.Reply("Unable to fetch entity")

		return
	}

	log.Debug(entity)

	DisplayEntity(ctx, entity, "Fetched data")
}
