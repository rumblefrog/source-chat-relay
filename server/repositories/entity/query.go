package entity

import (
	"database/sql"

	"github.com/rumblefrog/source-chat-relay/server/database"
	log "github.com/sirupsen/logrus"
)

func FetchEntity(id string, eType EntityType) (*Entity, error) {
	row := database.DBConnection.QueryRow("SELECT * FROM `relay_entities` WHERE `id` = ? AND `type` = ?", id, eType)

	var (
		entity          = &Entity{}
		receiveChannels string
		sendChannels    string
	)

	err := row.Scan(
		&entity.ID,
		&entity.Type,
		&receiveChannels,
		&sendChannels,
		&entity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	entity.ReceiveChannels = ParseChannels(receiveChannels)

	entity.SendChannels = ParseChannels(sendChannels)

	return entity, nil
}

func (entity *Entity) UpdateChannels() (sql.Result, error) {
	return database.DBConnection.Exec(
		"UPDATE `relay_entities` SET `receive_channels` = ?, `send_channels` = ? WHERE `id` = ?",
		EncodeChannels(entity.ReceiveChannels),
		EncodeChannels(entity.SendChannels),
		entity.ID,
	)
}

func (entity *Entity) CreateEntity() (sql.Result, error) {
	return database.DBConnection.Exec(
		"INSERT INTO `relay_entities` (`id`, `type`, `receive_channels`, `send_channels`) VALUES (?, ?, ?, ?)",
		entity.ID,
		entity.Type,
		EncodeChannels(entity.ReceiveChannels),
		EncodeChannels(entity.SendChannels),
	)
}

func FetchEntities(eType EntityType) ([]*Entity, error) {
	rows, err := database.DBConnection.Query("SELECT * FROM `relay_entities` WHERE `type` != ?", eType.Polarize())

	if err != nil {
		return nil, err
	}

	var entities []*Entity

	for rows.Next() {
		var (
			entity          = &Entity{}
			receiveChannels string
			sendChannels    string
		)

		if err = rows.Scan(
			&entity.ID,
			&entity.Type,
			&receiveChannels,
			&sendChannels,
			&entity.CreatedAt,
		); err != nil {
			rows.Close()
			return nil, err
		}

		entity.ReceiveChannels = ParseChannels(receiveChannels)

		entity.SendChannels = ParseChannels(sendChannels)

		entities = append(entities, entity)
	}

	return entities, nil
}

func CreateTable() {
	_, err := database.DBConnection.Exec("CREATE TABLE IF NOT EXISTS `relay_entities` ( `id` VARCHAR(64) NOT NULL , `type` TINYINT NOT NULL DEFAULT '0' , `receive_channels` VARCHAR(32) NOT NULL DEFAULT '' , `send_channels` VARCHAR(32) NOT NULL DEFAULT '' , `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP , PRIMARY KEY (`id`), INDEX (`type`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")

	if err != nil {
		log.WithField("error", err).Fatal("Unable to create tables")
	}
}
