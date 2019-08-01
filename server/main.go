package main

import (
	"github.com/kardianos/service"
	"github.com/rumblefrog/source-chat-relay/server/relay"
	"github.com/rumblefrog/source-chat-relay/server/ui"
	"github.com/sirupsen/logrus"

	"github.com/rumblefrog/source-chat-relay/server/entity"

	"github.com/rumblefrog/source-chat-relay/server/bot"
	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/database"
)

type program struct{}

func init() {
	logrus.SetLevel(logrus.InfoLevel)
}

func (p *program) Start(s service.Service) error {
	logrus.Infof("Server is now running on version %s. Press CTRL-C to exit.", config.SCRVER)

	config.ParseConfig()
	database.InitializeDatabase()
	entity.Initialize()

	relay.Instance = relay.NewRelay()
	relay.Instance.Listen(config.Config.General.Port)

	bot.Initialize()

	if config.Config.UI.Enabled {
		ui.UIListen()
	}

	return nil
}

func (p *program) Stop(s service.Service) error {
	logrus.Info("Received exit signal. Terminating.")

	bot.RelayBot.Close()
	relay.Instance.Listener.Close()
	database.Connection.Close()

	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "source-chat-relay",
		DisplayName: "Source Chat Relay",
		Description: "Service for Source Chat Relay",
	}

	prg := &program{}

	s, err := service.New(prg, svcConfig)

	if err != nil {
		logrus.Fatal(err)
	}

	err = s.Run()

	if err != nil {
		logrus.Fatal(err)
	}
}
