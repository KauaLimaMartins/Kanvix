package domain

import "strings"

func HasRole(role string, allowed ...string) bool {
	role = strings.ToLower(strings.TrimSpace(role))
	for _, a := range allowed {
		if role == strings.ToLower(strings.TrimSpace(a)) {
			return true
		}
	}
	return false
}

