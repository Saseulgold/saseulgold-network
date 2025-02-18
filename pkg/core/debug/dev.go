package debug

import (
	"fmt"
	C "hello/pkg/core/config"
)

func DebugLog(args ...interface{}) {
	// if C.CORE_TEST_MODE {
	// fmt.Printf("%s\n", args)
	//}
}

func DebugPanic(format string, args ...interface{}) {
	if C.CORE_TEST_MODE {
		panic(fmt.Sprintf(format, args...))
	}
}

func DebugAssert(condition bool, format string, args ...interface{}) {
	if C.CORE_TEST_MODE && !condition {
		panic(fmt.Sprintf(format, args...))
	}
}
