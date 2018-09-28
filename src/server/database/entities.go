package database

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type EntityType int

const (
	Server EntityType = iota
	Channel
)

type Entity struct {
	ID              string
	Type            EntityType
	ReceiveChannels []int
	SendChannels    []int
	CreatedAt       time.Time
}

func FetchEntity(id string) (*Entity, error) {
	stmt, err := DBConnection.Prepare("SELECT * FROM `relay_entities` WHERE `id` = ?")

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
	return DBConnection.Exec(
		"UPDATE `relay_entities` SET `receive_channels` = ?, `send_channels` = ?",
		EncodeChannels(entity.ReceiveChannels),
		EncodeChannels(entity.SendChannels),
	)
}

func ParseChannels(s string) (c []int) {
	ss := strings.Split(s, ",")

	for _, channel := range ss {
		tc, _ := strconv.Atoi(channel)
		c = append(c, tc)
	}

	return
}

func EncodeChannels(channels []int) string {
	var s []string

	for c := range channels {
		s = append(s, strconv.Itoa(c))
	}

	return strings.Join(s, ",")
}
