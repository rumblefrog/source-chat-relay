package entity

import (
	"strconv"
	"strings"

	"github.com/rumblefrog/source-chat-relay/server/protocol"
)

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
