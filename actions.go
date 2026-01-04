package main

import (
	"os"
	"os/exec"
)

type Action func()

// ---- Example Actions ----
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