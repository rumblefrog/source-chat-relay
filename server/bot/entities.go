package bot

import (
	"time"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/Necroforger/dgwidgets"

	repoEntity "github.com/rumblefrog/source-chat-relay/server/repositories/entity"
)

func EntitiesCMD(ctx *exrouter.Context) {
	eType := repoEntity.EntityTypeFromString(ctx.Args.Get(1))

	entities := repoEntity.GetEntities(eType)

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
