package bot

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/rumblefrog/source-chat-relay/server/entity"
)

func deleteCommand(ctx *exrouter.Context) {
	id := ctx.Args.Get(1)

	if len(id) == 0 {
		ctx.Reply("Missing ID")

		return
	}

	tEntity, err := entity.GetEntity(id)

	if err != nil {
		ctx.Reply("Unable to fetch entity")

		return
	}

	tEntity.Delete()

	ctx.Reply("Entity deleted")
}
