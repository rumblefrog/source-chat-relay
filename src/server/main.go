package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rumblefrog/source-chat-relay/src/server/database"
	"github.com/rumblefrog/source-chat-relay/src/server/helper"
	"github.com/rumblefrog/source-chat-relay/src/server/protocol"
)

func main() {
	helper.LoadConfig()

	database.InitDB()

	protocol.InitSocket()

	// bot.InitBot()

	log.Println("Server is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Received exit signal. Terminating.")

	// bot.RelayBot.Close()

	protocol.NetListener.Close()

	database.DBConnection.Close()
}
