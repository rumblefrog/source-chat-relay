package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/rumblefrog/source-chat-relay/src/server/bot"
	"github.com/rumblefrog/source-chat-relay/src/server/database"
	"github.com/rumblefrog/source-chat-relay/src/server/protocol"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Println("Server is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Received exit signal. Terminating.")

	bot.RelayBot.Session.Close()

	protocol.NetListener.Close()

	database.DBConnection.Close()
}
