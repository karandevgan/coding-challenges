package internal

import "strings"

func GetStringAtColumn(s string, sep string, col int, cols ...int) string {
	l := strings.Split(s, sep)
	ss := make([]string, 0, len(cols))
	ss = append(ss, getStringAtColumn(l, len(l), col))
	for _, c := range cols {
		ss = append(ss, getStringAtColumn(l, len(l), c))
	}
	return strings.Join(ss, sep)
}

func getStringAtColumn(l []string, sLen int, col int) string {
	if col <= 0 || col > sLen {
		return ""
	}
	return l[col-1]
}
