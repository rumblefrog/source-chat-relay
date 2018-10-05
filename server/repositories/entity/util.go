package entity

import (
	"strconv"
	"strings"
)

func (entity *Entity) CanReceive(channels []int) bool {
	for _, c := range entity.ReceiveChannels {
		for _, c1 := range channels {
			if c == c1 || (c == -1 && c1 != 0) || (c1 == -1 && c != 0) {
				return true
			}
		}
	}

	return false
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
		if c == 0 {
			continue
		}

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
