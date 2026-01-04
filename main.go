package main

import "fmt"

const vlcPath = "/usr/bin/vlc"

//
// ---- Main ----
//
func main() {
	jobs := []Job{
		{
			Name: "Morning Music",
			Time: "07:00",
			Run:  playMusic("/home/youruser/music/morning.mp4"),
		},
		{
			Name: "Cleanup Temp",
			Time: "02:00",
			Run:  runCommand("/usr/bin/rm", "-rf", "/tmp/myapp"),
		},
	}

	fmt.Println("Custom cron scheduler started")

	for _, job := range jobs {
		go runJob(job)
	}

	select {} // block forever
}
