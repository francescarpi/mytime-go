package logger

import (
	"log"
	"os"
)

func Init() *os.File {
	logFile := "mytime.log"

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Error opening or creating log file:")
	}

	log.SetOutput(file)

	return file
}
