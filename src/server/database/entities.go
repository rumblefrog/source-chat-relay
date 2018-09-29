package database

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type EntityType int

const (
	Server EntityType = iota
	Channel
	All
)

type Entity struct {
	ID              string
	Type            EntityType
	ReceiveChannels []int
	SendChannels    []int
	CreatedAt       time.Time
}

func FetchEntity(id string, eType EntityType) (*Entity, error) {
	row := DBConnection.QueryRow("SELECT * FROM `relay_entities` WHERE `id` = ? AND `type` = ?", id, eType)

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

func FetchEntities(eType EntityType) ([]*Entity, error) {
	rows, err := DBConnection.Query("SELECT * FROM `relay_entities` WHERE `type` != ?", eType.Polarize())

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

func (entity *Entity) UpdateChannels() (sql.Result, error) {
	return DBConnection.Exec(
		"UPDATE `relay_entities` SET `receive_channels` = ?, `send_channels` = ? WHERE `id` = ?",
		EncodeChannels(entity.ReceiveChannels),
		EncodeChannels(entity.SendChannels),
		entity.ID,
	)
}

func (entity *Entity) CreateEntity() (sql.Result, error) {
	return DBConnection.Exec(
		"INSERT INTO `relay_entities` (`id`, `type`, `receive_channels`, `send_channels`) VALUES (?, ?, ?, ?)",
		entity.ID,
		entity.Type,
		EncodeChannels(entity.ReceiveChannels),
		EncodeChannels(entity.SendChannels),
	)
}

func (entity *Entity) CanReceive(channels []int) bool {
	for _, c := range entity.ReceiveChannels {
		for _, c1 := range channels {
			if c == c1 || c == -1 {
				return true
			}
		}
	}

	return false
}

func (entity *Entity) Embed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: 0xE1C15C,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  entity.DisplayName(),
				Value: entity.ID,
			},
			&discordgo.MessageEmbedField{
				Name:  "Entity Type :gear:",
				Value: entity.Type.String(),
			},
			&discordgo.MessageEmbedField{
				Name:  "Receive Channels :inbox_tray:",
				Value: ChannelString(entity.ReceiveChannels),
			},
			&discordgo.MessageEmbedField{
				Name:  "Send Channels :outbox_tray:",
				Value: ChannelString(entity.SendChannels),
			},
			&discordgo.MessageEmbedField{
				Name:  "Created At",
				Value: entity.CreatedAt.Format(time.RFC1123),
			},
		},
	}
}

func (entity *Entity) DisplayName() string {
	if entity.Type == Server {
		return "Entity ID (Keep Private) :key:"
	}

	return "Entity ID :key:"
}

func ParseChannels(s string) (c []int) {
	ss := strings.Split(strings.Replace(s, " ", "", -1), ",")

	for _, channel := range ss {
		tc, _ := strconv.Atoi(channel)
		c = append(c, tc)
	}

	return
}

func EncodeChannels(channels []int) string {
	var s []string

	for _, c := range channels {
		s = append(s, strconv.Itoa(c))
	}

	return strings.Join(s, ",")
}

func ChannelString(channels []int) string {
	var s []string

	for _, c := range channels {
		s = append(s, strconv.Itoa(c))
	}

	j := strings.Join(s, ", ")

	if j == "" {
		return "None"
	}

	return j
}

func (eType EntityType) Polarize() EntityType {
	switch eType {
	case Server:
		return Channel
	case Channel:
		return Server
	default:
		return All
	}
}

func EntityTypeFromString(t string) EntityType {
	switch strings.ToLower(t) {
	case "server":
		return Server
	case "channel":
		return Channel
	default:
		return All
	}
}

func (eType EntityType) String() string {
	switch eType {
	case Server:
		return "Server"
	case Channel:
		return "Channel"
	default:
		return "Unknown"
	}
}
