package database

import (
	"database/sql"
	"log"
	"time"
)

type RelayEntities struct {
	ID              int
	Source          string
	Type            int
	ReceiveChannels string
	SendChannels    string
	CreatedAt       time.Time
}

func CreateTables(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `relay_entities` ( `id` INT NOT NULL AUTO_INCREMENT , `source` VARCHAR(64) NOT NULL , `type` TINYINT NOT NULL DEFAULT '0' , `receive_channels` VARCHAR(32) DEFAULT NULL , `send_channels` VARCHAR(32) DEFAULT NULL , `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP , PRIMARY KEY (`id`), INDEX (`source`), INDEX (`type`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")

	if err != nil {
		log.Panic("Unable to create tables", err)
	}
}
