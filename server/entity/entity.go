package entity

import "time"

type Entity struct {
	ID                   string
	DisplayName          string
	ReceiveChannels      []int
	SendChannels         []int
	DisabledReceiveTypes []int
	DisabledSendTypes    []int
	CreatedAt            time.Time
}

func Initialize() {
	initializeTable()
	upgradeSchema()
	loadEntitiesIntoCache()
}
