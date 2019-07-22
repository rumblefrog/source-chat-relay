package config

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var Conf Config

func ParseConfig() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		logrus.WithField("error", err).Fatal("Unable to parse config")
	}
}
