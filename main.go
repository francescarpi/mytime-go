package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mytime/config"
	"mytime/db"
	"mytime/interfaces/gui"
	"mytime/logger"
)

func main() {
	guiFlag := flag.Bool("gui", false, "Run the GUI")
	logsOn := flag.Bool("logs", false, "Enable logs")
	flag.Parse()

	if *logsOn {
		logFile := logger.Init()
		defer logFile.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	cfg := config.Config{}
	cfg.AutoDiscover()

	conn := db.GetConnection(&cfg)

	if *guiFlag {
		gui.Init(conn)
		return
	}

	fmt.Println("Usage: ")
	fmt.Println("  mytime [--gui]")

}
