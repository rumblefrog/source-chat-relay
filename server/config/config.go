package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
)

var Conf Config

func init() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		log.WithField("error", err).Fatal("Unable to parse config")
	}
}
