package bot

import (
	"github.com/Necroforger/dgrouter/exrouter"
	log "github.com/sirupsen/logrus"
)

func SetChannel(ctx *exrouter.Context) {
	log.Debug(ctx.Args)
}
