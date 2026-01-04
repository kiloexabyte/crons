package main

import (
	"fmt"
	"time"
)

//
// ---- Core Scheduler Types ----
//
type Job struct {
	Name string
	Time string // "HH:MM"
	Run  Action
}

//
// ---- Scheduler Logic ----
//
func nextRun(timeStr string) time.Time {
	now := time.Now()

	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		panic("Invalid time format: " + timeStr)
	}

	next := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		t.Hour(),
		t.Minute(),
		0,
		0,
		now.Location(),
	)

	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}

	return next
}

func runJob(job Job) {
	for {
		runAt := nextRun(job.Time)
		fmt.Printf("Scheduled [%s] at %s\n",
			job.Name,
			runAt.Format(time.RFC1123),
		)

		time.Sleep(time.Until(runAt))

		fmt.Printf("Running [%s]\n", job.Name)
		job.Run()

		time.Sleep(1 * time.Second)
	}
}