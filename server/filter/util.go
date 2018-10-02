package filter

import (
	"github.com/rumblefrog/source-chat-relay/server/helper"
)

func IsInFilter(s string) bool {
	if !helper.Conf.General.Filter {
		return false
	}

	if len(Filter) <= 0 {
		return false
	}

	for _, r := range Filter {
		if r.MatchString(s) {
			return true
		}
	}

	return false
}
