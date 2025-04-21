package main

import (
	"flag"
	"fmt"
	"mytime/config"
	"mytime/db"
	"mytime/interfaces/gui"
	"mytime/logger"
)

func main() {
	guiFlag := flag.Bool("gui", false, "Run the GUI")
	flag.Parse()

	logFile := logger.Init()
	defer logFile.Close()

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
