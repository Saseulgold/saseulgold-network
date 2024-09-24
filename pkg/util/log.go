package util

import (
	"log"
	"os"
)

var appLogger *log.Logger

func GetLogger() *log.Logger{
	if(appLogger== nil) {
			appLogger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	}

	return appLogger
}

