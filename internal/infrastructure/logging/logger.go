package logging

import (
	"fmt"
	"log"
)

var logger = log.Default()

func LogErrorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Error]: %s\n", msg)
}

func LogInfof(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}

func LogPanicf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Panicf("%s\n", msg)
}
