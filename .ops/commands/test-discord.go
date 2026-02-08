package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (Ops) Testdiscord() error {
	_ = godotenv.Load(".env")

	webhookURL := os.Getenv("DISCORD_WEBHOOK")
	if webhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK not set")
	}

	payload := map[string]string{"content": "🧪 **Test alert** - Discord webhook is working!"}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	fmt.Println("Discord test message sent!")
	return nil
}
