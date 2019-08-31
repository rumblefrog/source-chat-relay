package bot

import (
	"time"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/Necroforger/dgwidgets"
	"github.com/rumblefrog/source-chat-relay/server/entity"
)

func entitiesCommand(ctx *exrouter.Context) {
	entities := entity.Entities()

	if len(entities) <= 0 {
		ctx.Reply("No entities found in database")
	}

	p := dgwidgets.NewPaginator(ctx.Ses, ctx.Msg.ChannelID)

	for _, entity := range entities {
		p.Add(entity.Embed())
	}

	p.SetPageFooters()

	p.ColourWhenDone = 0x1DB954

	p.Widget.Timeout = time.Minute * 2

	p.Spawn()
}
