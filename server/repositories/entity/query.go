package entity

import (
	"database/sql"

	"github.com/rumblefrog/source-chat-relay/server/database"
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

	defer rows.Close()

	for rows.Next() {
		var (
			entity          = &Entity{}
			receiveChannels string
			sendChannels    string
		)

		rows.Scan(
			&entity.ID,
			&entity.Type,
			&receiveChannels,
			&sendChannels,
			&entity.CreatedAt,
		)

		entity.ReceiveChannels = ParseChannels(receiveChannels)

		entity.SendChannels = ParseChannels(sendChannels)

		entities = append(entities, entity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}
