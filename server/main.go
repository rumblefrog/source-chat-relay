package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rumblefrog/source-chat-relay/server/relay"
	"github.com/sirupsen/logrus"

	"github.com/rumblefrog/source-chat-relay/server/entity"

	"github.com/rumblefrog/source-chat-relay/server/bot"
	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/database"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	logrus.Infof("Server is now running on version %s. Press CTRL-C to exit.", config.SCRVER)

	config.ParseConfig()
	database.InitializeDatabase()
	entity.Initialize()
	bot.Initialize()

	relay.Instance = relay.NewRelay()
	relay.Instance.Listen(config.Conf.General.Port)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	logrus.Info("Received exit signal. Terminating.")

	bot.RelayBot.Close()
	relay.Instance.Listener.Close()
	database.Connection.Close()
}
