package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rumblefrog/source-chat-relay/src/server/helper"
	"github.com/rumblefrog/source-chat-relay/src/server/socket"
	"github.com/rumblefrog/source-chat-relay/src/server/storage"
)

func main() {
	helper.LoadConfig()

	storage.InitDB()

	socket.InitSocket()

	// bot.InitBot()

	log.Println("Server is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Received exit signal. Terminating.")

	// bot.Session.Close()

	socket.NetListener.Close()
}
