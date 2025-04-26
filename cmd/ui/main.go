package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/francescarpi/mytime/internal/ui"
)

func main() {
	logsOn := flag.Bool("logs", false, "Enable logs to file")
	flag.Parse()

	logFile := setupLogging(*logsOn)
	if logFile != nil {
		defer logFile.Close()
	}

	ui.StartApp()
}

func setupLogging(enabled bool) *os.File {
	if !enabled {
		log.SetOutput(io.Discard)
		return nil
	}

	logFilePath := "mytime.log"
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		os.Exit(1)
	}

	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("\n===== Logging started =====")

	return file
}

