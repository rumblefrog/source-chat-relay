package database

import (
	log "github.com/sirupsen/logrus"
)

type EntityCache struct {
	Entities   []*Entity
	Controller chan *Entity
}

var Cache *EntityCache

func (c *EntityCache) StartCache() {
	for {
		select {
		case entity := <-c.Controller:
			item := c.InCache(entity)

			if item != nil {
				item.ReceiveChannels = entity.ReceiveChannels
				item.SendChannels = entity.SendChannels
			} else {
				c.Entities = append(c.Entities, entity)
			}
		}
	}
}

func (c *EntityCache) InCache(entity *Entity) *Entity {
	for _, e := range c.Entities {
		if e.ID == entity.ID {
			return e
		}
	}

	return nil
}

func (c *EntityCache) GetEntity(id string) *Entity {
	for _, e := range c.Entities {
		if e.ID == id {
			return e
		}
	}

	return nil
}

func (c *EntityCache) DownloadCache() {
	entities, err := FetchEntities(All)

	if err != nil {
		log.WithField("error", err).Warn("Unable to sync relay channel cache")
		return
	}

	c.Entities = entities

	log.WithField("len", len(c.Entities)).Info("Relay channel cache synced")
}
