package util

import (
	"bytes"
	"strconv"
	"strings"
)

func PadLeft(str string, padChar string, length int) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(padChar, length-len(str)) + str
}

func PadRight(str string, padChar string, length int) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(padChar, length-len(str))
}

func Concat(s ...string) string {
	var buffer bytes.Buffer

	for i := 0; i < len(s); i++ {
		buffer.WriteString(s[i])
	}
	return buffer.String()
}

func String(value interface{}) string {
	switch v := value.(type) {
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		return v
	case []byte:
		return string(v)
	}

	return ""
}
