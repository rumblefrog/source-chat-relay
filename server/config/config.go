package config

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var Config Config_t

func ParseConfig() {
	if _, err := toml.DecodeFile("config.toml", &Config); err != nil {
		logrus.WithField("error", err).Fatal("Unable to parse config")
	}
}
