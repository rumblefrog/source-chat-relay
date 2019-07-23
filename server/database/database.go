package database

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/go-sql-driver/mysql"
	"github.com/rumblefrog/source-chat-relay/server/config"
)

var Connection *sql.DB

func InitializeDatabase() {
	c := mysql.NewConfig()

	c.Net = config.Config.Database.Protocol

	c.User = config.Config.Database.Username

	c.Passwd = config.Config.Database.Password

	c.DBName = config.Config.Database.Database

	c.Collation = "utf8mb4_general_ci"

	c.InterpolateParams = true

	c.ParseTime = true

	if config.Config.Database.Protocol == "tcp" {
		c.Addr = fmt.Sprintf("%s:%d", config.Config.Database.Host, config.Config.Database.Port)
	} else {
		c.Addr = config.Config.Database.Host
	}

	var err error

	Connection, err = sql.Open("mysql", c.FormatDSN())

	if err != nil {
		logrus.WithField("error", err).Fatal("Unable to connect to database")
	}
}
