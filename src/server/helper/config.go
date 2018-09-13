package helper

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Token string `toml:"Token"`
	Port  int    `toml:"Port"`
}

var Conf Config

func LoadConfig() *Config {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		log.Panic(err)
	}

	return &Conf
}
