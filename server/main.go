package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/rumblefrog/source-chat-relay/server/bot"
	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/database"
	"github.com/rumblefrog/source-chat-relay/server/protocol"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Infof("Server is now running on version %s. Press CTRL-C to exit.", config.SCRVER)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Info("Received exit signal. Terminating.")

	bot.RelayBot.Session.Close()

	protocol.NetListener.Close()

	database.DBConnection.Close()
}
