package git

import "strings"

func isSet(s ...string) bool {
	return (len(s) > 0) && (strings.TrimSpace(s[0]) != "")
}

func setOrDefault(def string, s ...string) string {
	if (len(s) == 0) || (strings.TrimSpace(s[0]) == "") {
		s = []string{def}
	}

	return strings.TrimSpace(s[0])
}
