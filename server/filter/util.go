package filter

import "strings"

func IsInFilter(s string) bool {
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
