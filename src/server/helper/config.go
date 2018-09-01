package helper

import (
	"log"

	"github.com/BurntSushi/toml"
)

// Config - Data scheme for config file
type Config struct {
	Token string `toml:"Token"`
	Port  int    `toml:"Port"`
}

// Conf - Loaded config
var Conf Config

// LoadConfig - Load and parse config
func LoadConfig() *Config {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		log.Panic(err)
	}

	return &Conf
}
