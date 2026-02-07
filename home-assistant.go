package main

import "crons/pkg/homeassistant"

func setAllLightsBrightness(percent int) Action {
	return func() {
		client := homeassistant.NewClient()
		_ = client.SetAllLightsBrightness(percent)
	}
}
