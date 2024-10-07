package util

import (
	"strings"
)

func PadLeft(str string, padChar string, length int) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(padChar, length-len(str)) + str
}

