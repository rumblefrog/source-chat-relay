package filter

import "github.com/rumblefrog/source-chat-relay/server/config"

func IsInFilter(s string) bool {
	if !config.Config.General.Filter {
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
