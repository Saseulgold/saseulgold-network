package util

import (
	"time"
)

func Time() int64 {
	return time.Now().Unix()
}

/**
func UTime() int64 {
	microtime := time.Now().UnixNano()
	return ((microtime / 1000) / 1000) * 1000
}
*/
