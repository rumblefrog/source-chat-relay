package entity

import (
	log "github.com/sirupsen/logrus"
)

var Cache *EntityCache

func Start() {
	go Sync()

	for {
		select {
		case entity := <-Cache.Controller:
			Cache.Entities[entity.ID] = entity
		}
	}
}

func InCache(key string, cType EntityType) (string, bool) {
	for k, e := range Cache.Entities {
		if k == key && (e.Type == cType || cType == All) {
			return k, true
		}
	}

	return "", false
}

func Sync() {
	entities, err := FetchEntities(All)

	if err != nil {
		log.WithField("error", err).Warn("Unable to sync relay channel cache")
		return
	}

	for _, e := range entities {
		Cache.Entities[e.ID] = e
	}

	log.WithField("len", len(Cache.Entities)).Info("Relay channel cache synced")
}
