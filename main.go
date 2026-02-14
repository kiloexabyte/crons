package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

// ---- Main ----
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env: %v", err)
	}

	jobs := []Job{
		{
			Name: "Dim Lights",
			Time: "20:00",
			Run:  setAllLightsBrightness(10),
		},
		{
			Name: "Brighten Lights",
			Time: "05:00",
			Run:  setAllLightsBrightness(100),
		},
		{
			Name: "Heater On",
			Time: "05:00",
			Run:  setHeater(true),
		},
		{
			Name: "Heater Off",
			Time: "06:00",
			Run:  setHeater(false),
		},
	}

	intervalJobs := []IntervalJob{
		{
			Name:     "Health Check",
			Interval: 5 * time.Minute,
			Run:      checkServices(),
		},
	}

	fmt.Println("Custom cron scheduler started")

	for _, job := range jobs {
		go runJob(job)
	}

	for _, job := range intervalJobs {
		go runIntervalJob(job)
	}

	select {} // block forever
}
