package entity

import (
	"strconv"
	"strings"

	"github.com/rumblefrog/source-chat-relay/server/protocol"
)

func (e *Entity) ReceiveIntersectsWith(chans []int) bool {
	for _, e := range e.ReceiveChannels {
		for _, v := range chans {
			if e == 0 || v == 0 {
				continue
			}

			if e == v || e == -1 || v == -1 {
				return true
			}
		}
	}

	return false
}

func (e *Entity) SendIntersectsWith(chans []int) bool {
	for _, e := range e.SendChannels {
		for _, v := range chans {
			if e == 0 || v == 0 {
				continue
			}

			if e == v || e == -1 || v == -1 {
				return true
			}
		}
	}

	return false
}

func (e *Entity) CanReceiveType(t protocol.MessageType) bool {
	for _, v := range e.DisabledReceiveTypes {
		if protocol.MessageType(v) == t {
			return false
		}
	}

	return true
}

func (e *Entity) CanSendType(t protocol.MessageType) bool {
	for _, v := range e.DisabledSendTypes {
		if protocol.MessageType(v) == t {
			return false
		}
	}

	return true
}

func ParseDelimitedChannels(s string) (c []int) {
	ss := strings.Split(strings.Replace(s, " ", "", -1), ",")

	for _, channel := range ss {
		tc, _ := strconv.Atoi(channel)
		c = append(c, tc)
	}

	return
}

func EncodeDelimitedChannels(channels []int) string {
	var s []string

	for _, c := range channels {
		s = append(s, strconv.Itoa(c))
	}

	return strings.Join(s, ",")
}

func HumanizeChannelString(channels []int) string {
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

func DeliverableSendChannels(d protocol.Deliverable) []int {
	e, err := GetEntity(d.Author())

	if err != nil {
		return []int{}
	}

	return e.SendChannels
}
