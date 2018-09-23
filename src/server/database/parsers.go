package database

import (
	"strconv"
	"strings"
)

func ParseChannels(s string) (c []int) {
	ss := strings.Split(s, ",")

	for _, channel := range ss {
		tc, _ := strconv.Atoi(channel)
		c = append(c, tc)
	}

	return
}
