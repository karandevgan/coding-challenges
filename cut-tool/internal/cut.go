package internal

import "strings"

func GetStringAtColumn(s string, col int, sep string) string {
	l := strings.Split(s, sep)
	if col <= 0 {
		return ""
	}
	if len(l) < col {
		return ""
	}
	return l[col-1]
}
