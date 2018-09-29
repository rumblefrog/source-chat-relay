package bot

import (
	"time"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/Necroforger/dgwidgets"
	"github.com/rumblefrog/source-chat-relay/src/server/database"
	log "github.com/sirupsen/logrus"
)

func EntitiesCMD(ctx *exrouter.Context) {
	eType := database.EntityTypeFromString(ctx.Args.Get(1))

	entities, err := database.FetchEntities(eType)

	if err != nil {
		log.WithField("error", err).Warn("Unable to fetch entities")

		ctx.Reply("Unable to fetch entities")

		return
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
