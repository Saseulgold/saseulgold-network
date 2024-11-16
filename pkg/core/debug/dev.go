package debug

import (
	"fmt"
)

func DebugLog(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func DebugPanic(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
