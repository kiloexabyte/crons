package commands

import (
	"fmt"
	"os"
	"strings"

	"crons/pkg/homeassistant"

	"github.com/joho/godotenv"
)

func (Ops) Testheater() error {
	_ = godotenv.Load(".env")

	state := os.Getenv("STATE")
	if state == "" {
		return fmt.Errorf("usage: STATE=on op testheater (or STATE=off)")
	}

	on := strings.ToLower(state) == "on"

	client := homeassistant.NewClient()
	fmt.Printf("Setting heater to %s...\n", strings.ToUpper(state))
	if err := client.SetHeater(on); err != nil {
		return fmt.Errorf("failed to set heater: %w", err)
	}
	fmt.Println("Done!")
	return nil
}
