package entity

import "github.com/sirupsen/logrus"

type EntityCache map[string]*Entity

var Cache EntityCache

func init() {
	Cache = make(map[string]*Entity)
}

func WriteCache(e *Entity) {
	Cache[e.ID] = e
}

func GetEntity(id string) (*Entity, error) {
	val, ok := Cache[id]

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
	for _, e := range Cache {
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
