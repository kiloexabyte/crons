package main

import (
	"fmt"

	"crons/pkg/homeassistant"
)

func setAllLightsBrightness(percent int) Action {
	return func() {
		client := homeassistant.NewClient()
		_ = client.SetAllLightsBrightness(percent)
	}
}

func setHeater(on bool) Action {
	return func() {
		client := homeassistant.NewClient()
		state := "OFF"
		if on {
			state = "ON"
		}
		if err := client.SetHeater(on); err != nil {
			fmt.Printf("Error setting heater: %v\n", err)
			return
		}
		fmt.Printf("Heater set to %s\n", state)
	}
}
