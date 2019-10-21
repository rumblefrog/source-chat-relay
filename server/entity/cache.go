package entity

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type EntityCache struct {
	sync.RWMutex

	Entities map[string]*Entity
}

var Cache *EntityCache

func init() {
	Cache = &EntityCache{
		Entities: make(map[string]*Entity),
	}
}

func WriteCache(e *Entity) {
	Cache.Lock()
	defer Cache.Unlock()

	Cache.Entities[e.ID] = e
}

func GetEntity(id string) (*Entity, error) {
	Cache.RLock()
	defer Cache.RUnlock()

	val, ok := Cache.Entities[id]

	if ok {
		return val, nil
	}

	entity, err := FetchEntity(id)

	if err != nil {
		return nil, err
	}

	WriteCache(entity)

	return entity, nil
}

func Entities() (entities []*Entity) {
	Cache.RLock()
	defer Cache.RUnlock()

	for _, e := range Cache.Entities {
		entities = append(entities, e)
	}

	return
}

func loadEntitiesIntoCache() {
	entities, err := FetchEntities()

	if err != nil {
		logrus.WithField("error", err).Warn("Unable to initialize entities cache")

		return
	}

	for _, e := range entities {
		WriteCache(e)
	}

	logrus.WithField("length", len(entities)).Info("Entities cache initialized")
}
