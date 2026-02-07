package main

import (
	"fmt"
	"log"

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
	}

	fmt.Println("Custom cron scheduler started")

	for _, job := range jobs {
		go runJob(job)
	}

	select {} // block forever
}
