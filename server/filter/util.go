package filter

import (
	"strings"

	"github.com/rumblefrog/source-chat-relay/server/helper"
)

func IsInFilter(s string) bool {
	if !helper.Conf.General.Filter {
		return false
	}

	if len(Filter) <= 0 {
		return false
	}

	for _, i := range Filter {
		if strings.Contains(s, i) {
			return true
		}
	}

	return false
}
