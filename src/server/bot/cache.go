package bot

import (
	"github.com/rumblefrog/source-chat-relay/src/server/database"
	log "github.com/sirupsen/logrus"
)

func (b *DiscordBot) StartCache() {
	for {
		select {
		case entity := <-b.CacheController:
			index := b.GetCacheIndex(entity)

			if index > -1 {
				b.Cache[index] = entity
			} else {
				b.Cache = append(b.Cache, entity)
			}
		}
	}
}

func (b *DiscordBot) GetCacheIndex(entity *database.Entity) int {
	for i, e := range b.Cache {
		if e.ID == entity.ID {
			return i
		}
	}

	return -1
}

func (b *DiscordBot) SyncCache() {
	entities, err := database.FetchEntities(database.Channel)

	if err != nil {
		log.WithField("error", err).Warn("Unable to sync channel cache")
		return
	}

	b.Cache = entities
}
