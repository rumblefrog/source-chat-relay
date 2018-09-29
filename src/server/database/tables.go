package database

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

func CreateTables(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `relay_entities` ( `id` VARCHAR(64) NOT NULL , `type` TINYINT NOT NULL DEFAULT '0' , `receive_channels` VARCHAR(32) NOT NULL DEFAULT '' , `send_channels` VARCHAR(32) NOT NULL DEFAULT '' , `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP , PRIMARY KEY (`id`), INDEX (`type`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")

	if err != nil {
		log.WithField("error", err).Fatal("Unable to create tables")
	}
}
