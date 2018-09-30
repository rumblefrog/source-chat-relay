package database

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/go-sql-driver/mysql"
	"github.com/rumblefrog/source-chat-relay/server/helper"
)

var DBConnection *sql.DB

func init() {
	c := mysql.NewConfig()

	c.Net = helper.Conf.Database.Protocol

	c.User = helper.Conf.Database.Username

	c.Passwd = helper.Conf.Database.Password

	c.DBName = helper.Conf.Database.Database

	c.Collation = "utf8mb4_general_ci"

	c.ParseTime = true

	if helper.Conf.Database.Protocol == "tcp" {
		c.Addr = fmt.Sprintf("%s:%d", helper.Conf.Database.Host, helper.Conf.Database.Port)
	} else {
		c.Addr = helper.Conf.Database.Host
	}

	var err error

	DBConnection, err = sql.Open("mysql", c.FormatDSN())

	if err != nil {
		log.WithField("error", err).Fatal("Unable to connect to database")
	}

	CreateTables(DBConnection)
}
