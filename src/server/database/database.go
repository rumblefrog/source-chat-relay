package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/rumblefrog/source-chat-relay/src/server/helper"
)

var DBConnection *sql.DB

func InitDB() {
	c := mysql.NewConfig()

	c.Addr = fmt.Sprintf("%s:%d", helper.Conf.Database.Host, helper.Conf.Database.Port)

	c.User = helper.Conf.Database.Username

	c.Passwd = helper.Conf.Database.Password

	c.DBName = helper.Conf.Database.Database

	c.Collation = "utf8mb4_general_ci"

	c.ParseTime = true

	var err error

	DBConnection, err = sql.Open("mysql", c.FormatDSN())

	if err != nil {
		log.Panic("Unable to connect to database")
	}

	CreateTables(DBConnection)
}
