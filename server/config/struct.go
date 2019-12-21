package config

type GeneralConfig struct {
	Port   int  `toml:"Port"`
	Filter bool `toml:"Filter"`
}

type BotConfig struct {
	Token         string `toml:"Token"`
	SimpleMessage bool   `toml:"SimpleMessage"`
	ListenToBots  bool   `toml:"ListenToBots"`
}

type DatabaseConfig struct {
	Host     string `toml:"Host"`
	Username string `toml:"Username"`
	Password string `toml:"Password"`
	Database string `toml:"Database"`
	Port     int    `toml:"Port"`
	Protocol string `toml:"Protocol"`
}

type UIConfig struct {
	Enabled bool `toml:"Enabled"`
	Port    int  `toml:"Port"`
}

type MsgConfig struct {
	EventFormatSimple                 string `toml:"EventFormatSimple"`
	EventFormatSimpleMapStart         string `toml:"EventFormatSimpleMapStart"`
	EventFormatSimpleMapEnd           string `toml:"EventFormatSimpleMapEnd"`
	EventFormatSimplePlayerConnect    string `toml:"EventFormatSimplePlayerConnect"`
	EventFormatSimplePlayerDisconnect string `toml:"EventFormatSimplePlayerDisconnect"`
	EventFormatSimplePlayerChat       string `toml:"EventFormatSimplePlayerChat"`
}

type Config_t struct {
	General  GeneralConfig
	Bot      BotConfig
	Database DatabaseConfig
	UI       UIConfig
	Messages MsgConfig
}
