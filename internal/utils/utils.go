package utils

import "strings"

func MixCase(s string) string {
	runes := []rune(strings.ToLower(s))
	for i := 0; i < len(runes); i += 2 {
		upper := []rune(strings.ToUpper(string(runes[i])))
		if len(upper) > 0 {
			runes[i] = upper[0]
		}
	}
	return string(runes)
}
