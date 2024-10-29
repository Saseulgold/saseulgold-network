package util

import (
	"bytes"
	"strings"
)

func PadLeft(str string, padChar string, length int) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(padChar, length-len(str)) + str
}

func Concat(s ...string) string {
	var buffer bytes.Buffer

	for i := 0; i < len(s); i++ {
		buffer.WriteString(s[i])
	}
	return buffer.String()
}
