package cron

import (
	"log"
	"time"

	"github.com/TejasGhatte/fampay-task-2024/routines"
	"github.com/robfig/cron/v3"
)

var Scheduler *cron.Cron

func Init() {
	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Println("Error loading location for crons: " + err.Error())
	}
	Scheduler = cron.New(cron.WithSeconds(), cron.WithLocation(location))

	_, err = Scheduler.AddFunc("*/10 * * * * *", func() {
		go routines.FetchVideos()
	})

	if err != nil {
		log.Println("Error adding cron job: " + err.Error())
	}

	Scheduler.Start()
}