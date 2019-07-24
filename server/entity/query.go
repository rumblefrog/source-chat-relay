package entity

import (
	"database/sql"

	"github.com/rumblefrog/source-chat-relay/server/database"
	"github.com/sirupsen/logrus"
)

func FetchEntity(id string) (*Entity, error) {
	// Specify the column names for backward compat with old database
	stmt, err := database.Connection.Prepare("SELECT `id`, `display_name`, `receive_channels`, `send_channels`, `disabled_receive_types`, `disabled_send_types`, `created_at` FROM `relay_entities` WHERE `id` = ?")

	if err != nil {
		return nil, err
	}

	var (
		entity               = &Entity{}
		receiveChannels      string
		sendChannels         string
		disabledReceiveTypes string
		disabledSendTypes    string
	)

	err = stmt.QueryRow(id).Scan(
		&entity.ID,
		&entity.DisplayName,
		&receiveChannels,
		&sendChannels,
		&disabledReceiveTypes,
		&disabledSendTypes,
		&entity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	entity.ReceiveChannels = ParseDelimitedChannels(receiveChannels)
	entity.SendChannels = ParseDelimitedChannels(sendChannels)
	entity.DisabledReceiveTypes = ParseDelimitedChannels(disabledReceiveTypes)
	entity.DisabledSendTypes = ParseDelimitedChannels(disabledSendTypes)

	return entity, nil
}

func FetchEntities() ([]*Entity, error) {
	// Specify the column names for backward compat with old database
	rows, err := database.Connection.Query("SELECT `id`, `display_name`, `receive_channels`, `send_channels`, `disabled_receive_types`, `disabled_send_types`, `created_at` FROM `relay_entities`")

	if err != nil {
		return nil, err
	}

	var entities []*Entity

	for rows.Next() {
		var (
			entity               = &Entity{}
			receiveChannels      string
			sendChannels         string
			disabledReceiveTypes string
			disabledSendTypes    string
		)

		if err = rows.Scan(
			&entity.ID,
			&entity.DisplayName,
			&receiveChannels,
			&sendChannels,
			&disabledReceiveTypes,
			&disabledSendTypes,
			&entity.CreatedAt,
		); err != nil {
			rows.Close()

			return nil, err
		}

		entity.ReceiveChannels = ParseDelimitedChannels(receiveChannels)
		entity.SendChannels = ParseDelimitedChannels(sendChannels)
		entity.DisabledReceiveTypes = ParseDelimitedChannels(disabledReceiveTypes)
		entity.DisabledSendTypes = ParseDelimitedChannels(disabledSendTypes)

		entities = append(entities, entity)
	}

	return entities, nil
}

func (entity *Entity) UpdateEntity() (sql.Result, error) {
	return database.Connection.Exec(
		"UPDATE `relay_entities` SET `display_name` = ?, `receive_channels` = ?, `send_channels` = ?, `disabled_receive_types` = ?, `disabled_send_types` = ? WHERE `id` = ?",
		entity.DisplayName,
		EncodeDelimitedChannels(entity.ReceiveChannels),
		EncodeDelimitedChannels(entity.SendChannels),
		EncodeDelimitedChannels(entity.DisabledReceiveTypes),
		EncodeDelimitedChannels(entity.DisabledSendTypes),
		entity.ID,
	)
}

func (entity *Entity) CreateEntity() (sql.Result, error) {
	return database.Connection.Exec(
		"INSERT INTO `relay_entities` (`id`, `display_name`, `receive_channels`, `send_channels`, `disabled_receive_types`, `disabled_send_types`) VALUES (?, ?, ?, ?, ?)",
		entity.ID,
		entity.DisplayName,
		EncodeDelimitedChannels(entity.ReceiveChannels),
		EncodeDelimitedChannels(entity.SendChannels),
		EncodeDelimitedChannels(entity.DisabledReceiveTypes),
		EncodeDelimitedChannels(entity.DisabledSendTypes),
	)
}

// Upgrading from v1 schema
func upgradeSchema() {
	// Ignore errors, if they fail, it doesn't exist or already does
	database.Connection.Exec("ALTER TABLE `relay_entities` DROP `type`")
	database.Connection.Exec("ALTER TABLE `relay_entities` ADD `disabled_receive_types` VARCHAR(32) NOT NULL DEFAULT '' AFTER `send_channels`")
	database.Connection.Exec("ALTER TABLE `relay_entities` ADD `disabled_send_types` VARCHAR(32) NOT NULL DEFAULT '' AFTER `disabled_receive_types`")
}

func initializeTable() {
	_, err := database.Connection.Exec("CREATE TABLE IF NOT EXISTS `relay_entities` ( `id` VARCHAR(64) NOT NULL , `display_name` VARCHAR(64) NOT NULL DEFAULT '' , `receive_channels` VARCHAR(32) NOT NULL DEFAULT '' , `send_channels` VARCHAR(32) NOT NULL DEFAULT '' , `disabled_receive_types` VARCHAR(32) NOT NULL DEFAULT '' , `disabled_send_types` VARCHAR(32) NOT NULL DEFAULT '' , `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP , PRIMARY KEY (`id`)) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE utf8mb4_general_ci;")

	if err != nil {
		logrus.WithField("error", err).Fatal("Unable to create tables")
	}
}
