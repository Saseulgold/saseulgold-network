package util

import (
	"fmt"
	t "time"
)

func CurrentTime(utime *int64) int64 {
	if utime != nil {
		return *utime / 1000000
	}
	return t.Now().Unix()
}

func Utime() int64 {
	return t.Now().UnixMicro()
}

func Uceiltime(utime *int64) int64 {
	if utime != nil {
		return Ufloortime(utime) + 1
	}
	return (t.Now().Unix() + 1) * 1000000
}

func Ufloortime(utime *int64) int64 {
	return CurrentTime(utime) * 1000000
}

func Bytetime(utime *int64) string {
	var t int64
	if utime != nil {
		t = *utime
	} else {
		t = Utime()
	}
	return fmt.Sprintf("%b", t)
}
