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

	c.Net = config.Conf.Database.Protocol

	c.User = config.Conf.Database.Username

	c.Passwd = config.Conf.Database.Password

	c.DBName = config.Conf.Database.Database

	c.Collation = "utf8mb4_general_ci"

	c.InterpolateParams = true

	c.ParseTime = true

	if config.Conf.Database.Protocol == "tcp" {
		c.Addr = fmt.Sprintf("%s:%d", config.Conf.Database.Host, config.Conf.Database.Port)
	} else {
		c.Addr = config.Conf.Database.Host
	}

	var err error

	Connection, err = sql.Open("mysql", c.FormatDSN())

	if err != nil {
		logrus.WithField("error", err).Fatal("Unable to connect to database")
	}
}
