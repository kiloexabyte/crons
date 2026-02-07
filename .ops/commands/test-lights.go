package commands

import (
	"fmt"
	"os"
	"strconv"

	"crons/pkg/homeassistant"

	"github.com/joho/godotenv"
)

func (Ops) Testlights() error {
	_ = godotenv.Load(".env")

	brightness := os.Getenv("BRIGHTNESS")
	if brightness == "" {
		return fmt.Errorf("usage: BRIGHTNESS=50 op testlights")
	}

	percent, err := strconv.Atoi(brightness)
	if err != nil {
		return fmt.Errorf("invalid brightness value: %s", brightness)
	}

	client := homeassistant.NewClient()
	fmt.Printf("Setting all lights to %d%%...\n", percent)
	return client.SetAllLightsBrightness(percent)
}
