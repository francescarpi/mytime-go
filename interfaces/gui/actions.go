package gui

import (
	"log"
	"time"
)

func GotoPreviousDate() {
	log.Println("Goto previous date. Current is ", date.Format(time.DateOnly))
	date = date.AddDate(0, 0, -1)
}

func GotoNextDate() bool {
	log.Println("Goto next date. Current is ", date.Format(time.DateOnly))
	nextDate := date.AddDate(0, 0, 1)
	if nextDate.After(time.Now()) {
		log.Println("Cannot go to future date")
		return false
	}
	date = nextDate
	return true
}
