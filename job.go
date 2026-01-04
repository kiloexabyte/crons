package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

//
// ---- Core Scheduler Types ----
//
type Action func()

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

//
// ---- Example Actions ----
//
func playMusic(file string) Action {
	return func() {
		cmd := exec.Command(
			vlcPath,
			"--intf", "dummy",
			"--play-and-exit",
			file,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}
}

func runCommand(command string, args ...string) Action {
	return func() {
		cmd := exec.Command(command, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}
}