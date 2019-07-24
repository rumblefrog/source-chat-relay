package entity

import "time"

type Entity struct {
	ID              string
	DisplayName     string
	ReceiveChannels []int
	SendChannels    []int
	CreatedAt       time.Time
}

func (e *Entity) ReceiveIntersectsWith(chans []int) bool {
	for _, e := range e.ReceiveChannels {
		for _, v := range chans {
			if e == v || (e == -1 && v != 0) || (v == -1 && e != 0) {
				return true
			}
		}
	}

	return false
}

func Initialize() {
	initializeTable()
	loadEntitiesIntoCache()
}
