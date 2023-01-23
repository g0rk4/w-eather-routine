package main

import (
	"log"
	"time"
	"weather-routine/executor"
	"weather-routine/scheduler"
)

func main() {
	sc := scheduler.NewScheduler(15 * time.Minute)
	ewExecutor := executor.EcoWatExecutor{}

	err := sc.Add(ewExecutor).Add(executor.EcoWatExecutor{}).Start(true)
	if err != nil {
		log.Fatal("Unable to start scheduler, program stopped. ", err)
	}
}
