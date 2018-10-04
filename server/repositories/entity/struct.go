package entity

import "time"

type EntityType int

const (
	Server EntityType = iota
	Channel
	All
)

type Entity struct {
	ID              string
	DisplayName     string
	Type            EntityType
	ReceiveChannels []int
	SendChannels    []int
	CreatedAt       time.Time
}

type EntityCache struct {
	Entities   map[string]*Entity
	Controller chan *Entity
}
