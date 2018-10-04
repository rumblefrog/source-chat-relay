package entity

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func GetEntity(id string, eType EntityType) (*Entity, error) {
	key, found := InCache(id, eType)

	if found {
		return Cache.Entities[key], nil
	}

	entity, err := FetchEntity(id, eType)

	if err != nil {
		return nil, err
	}

	Cache.Controller <- entity

	return entity, nil
}

func GetEntities(eType EntityType) (entities []*Entity) {
	for _, e := range Cache.Entities {
		if e.Type == eType || eType == All {
			entities = append(entities, e)
		}
	}

	return
}

func (entity *Entity) GetIDTitle() string {
	if entity.Type == Server {
		return "Entity ID (Keep Private) :key:"
	}

	return "Entity ID :key:"
}

func (entity *Entity) Embed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: 0xE1C15C,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  entity.GetIDTitle(),
				Value: entity.ID,
			},
			&discordgo.MessageEmbedField{
				Name:  "Display Name",
				Value: entity.DisplayName,
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
