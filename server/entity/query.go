package entity

import (
	"database/sql"

	"github.com/rumblefrog/source-chat-relay/server/database"
	"github.com/sirupsen/logrus"
)

func FetchEntity(id string) (*Entity, error) {
	// Specify the column names for backward compat with old database
	stmt, err := database.Connection.Prepare("SELECT `id`, `display_name`, `receive_channels`, `send_channels`, `created_at` FROM `relay_entities` WHERE `id` = ?")

	if err != nil {
		return nil, err
	}

	var (
		entity          = &Entity{}
		receiveChannels string
		sendChannels    string
	)

	err = stmt.QueryRow(id).Scan(
		&entity.ID,
		&entity.DisplayName,
		&receiveChannels,
		&sendChannels,
		&entity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	entity.ReceiveChannels = ParseDelimitedChannels(receiveChannels)
	entity.SendChannels = ParseDelimitedChannels(sendChannels)

	return entity, nil
}

func FetchEntities() ([]*Entity, error) {
	// Specify the column names for backward compat with old database
	rows, err := database.Connection.Query("SELECT `id`, `display_name`, `receive_channels`, `send_channels`, `created_at` FROM `relay_entities`")

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
			&entity.DisplayName,
			&receiveChannels,
			&sendChannels,
			&entity.CreatedAt,
		); err != nil {
			rows.Close()

			return nil, err
		}

		entity.ReceiveChannels = ParseDelimitedChannels(receiveChannels)
		entity.SendChannels = ParseDelimitedChannels(sendChannels)

		entities = append(entities, entity)
	}

	return entities, nil
}

func (entity *Entity) UpdateEntity() (sql.Result, error) {
	stmt, err := database.Connection.Prepare("UPDATE `relay_entities` SET `display_name` = ?, `receive_channels` = ?, `send_channels` = ? WHERE `id` = ?")

	if err != nil {
		return nil, err
	}

	return stmt.Exec(
		entity.DisplayName,
		EncodeDelimitedChannels(entity.ReceiveChannels),
		EncodeDelimitedChannels(entity.SendChannels),
		entity.ID,
	)
}

func (entity *Entity) CreateEntity() (sql.Result, error) {
	stmt, err := database.Connection.Prepare("INSERT INTO `relay_entities` (`id`, `display_name`, `receive_channels`, `send_channels`) VALUES (?, ?, ?, ?)")

	if err != nil {
		return nil, err
	}

	return stmt.Exec(
		entity.ID,
		entity.DisplayName,
		EncodeDelimitedChannels(entity.ReceiveChannels),
		EncodeDelimitedChannels(entity.SendChannels),
	)
}

func initializeTable() {
	_, err := database.Connection.Exec("CREATE TABLE IF NOT EXISTS `relay_entities` ( `id` VARCHAR(64) NOT NULL , `display_name` VARCHAR(64) NOT NULL DEFAULT '' , `receive_channels` VARCHAR(32) NOT NULL DEFAULT '' , `send_channels` VARCHAR(32) NOT NULL DEFAULT '' , `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP , PRIMARY KEY (`id`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")

	if err != nil {
		logrus.WithField("error", err).Fatal("Unable to create tables")
	}
}
