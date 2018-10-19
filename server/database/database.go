package database

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/go-sql-driver/mysql"
	"github.com/rumblefrog/source-chat-relay/server/config"
)

var DBConnection *sql.DB

func init() {
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

	DBConnection, err = sql.Open("mysql", c.FormatDSN())

	if err != nil {
		log.WithField("error", err).Fatal("Unable to connect to database")
	}
}
