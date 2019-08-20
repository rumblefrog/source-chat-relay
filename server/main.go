package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kardianos/service"
	"github.com/rumblefrog/source-chat-relay/server/relay"
	"github.com/rumblefrog/source-chat-relay/server/ui"
	"github.com/sirupsen/logrus"

	"github.com/rumblefrog/source-chat-relay/server/entity"

	"github.com/rumblefrog/source-chat-relay/server/bot"
	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/rumblefrog/source-chat-relay/server/database"
)

var (
	action string
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)

	flag.StringVar(&action, "service", "", "Install, uninstall, start, stop, restart")
	flag.StringVar(&config.Path, "config", "config.toml", "Path to the config file")

	flag.Parse()
}

type program struct{}

func (p *program) Start(s service.Service) error {
	logrus.Infof("Server is now running on version %s. Press CTRL-C to exit.", config.SCRVER)

	config.ParseConfig()
	database.InitializeDatabase()
	entity.Initialize()

	relay.Instance = relay.NewRelay()
	err := relay.Instance.Listen(config.Config.General.Port)

	if err != nil {
		logrus.WithField("error", err).Fatal("Unable to start listener")
	}

	bot.Initialize()

	if config.Config.UI.Enabled {
		go ui.UIListen()
	}

	return nil
}

func (p *program) Stop(s service.Service) error {
	logrus.Info("Received exit signal. Terminating.")

	bot.RelayBot.Close()

	relay.Instance.Closed = true
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

	flag.VisitAll(func(f *flag.Flag) {
		// ignore our own flags
		if f.Name == "service" {
			return
		}

		// ignore flags with default value
		if f.Value.String() == f.DefValue {
			return
		}

		svcConfig.Arguments = append(svcConfig.Arguments, "-"+f.Name+"="+f.Value.String())
	})

	s, err := service.New(&program{}, svcConfig)

	if err != nil {
		exit(err)
	}

	if action != "" {
		exit(actionHandler(action, s))
	}

	exit(s.Run())
}

func exit(err error) {
	if err != nil {
		logrus.Fatal(err)
	}

	os.Exit(0)
}

func actionHandler(action string, s service.Service) error {
	if action != "status" {
		return service.Control(s, action)
	}

	code, _ := s.Status()

	switch code {
	case service.StatusUnknown:
		fmt.Println("Service is not installed.")
	case service.StatusStopped:
		fmt.Println("Service is not running.")
	case service.StatusRunning:
		fmt.Println("Service is running.")
	default:
		fmt.Println("Error: ", code)
	}

	return nil
}
