package helper

import (
	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
)

type GeneralConfig struct {
	Port int `toml:"Port"`
}

type BotConfig struct {
	Token string `toml:"Token"`
}

type DatabaseConfig struct {
	Host     string `toml:"Host"`
	Username string `toml:"Username"`
	Password string `toml:"Password"`
	Database string `toml:"Database"`
	Port     int    `toml:"Port"`
}

type Config struct {
	General  GeneralConfig
	Bot      BotConfig
	Database DatabaseConfig
}

var Conf Config

func init() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		log.Fatal(err)
	}
}
