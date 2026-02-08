package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var nasIPs = []string{"192.168.50.11", "192.168.50.12"}

var services = []struct {
	Name string
	Port string
}{
	{"Jellyfin", "8096"},
	{"Sonarr", "8989"},
	{"Radarr", "8310"},
}

func checkServices() Action {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Track previous state for each service
	wasUp := make(map[string]bool)

	return func() {
		for _, svc := range services {
			up := false
			var lastErr error

			for _, ip := range nasIPs {
				url := fmt.Sprintf("http://%s:%s/", ip, svc.Port)
				resp, err := client.Get(url)
				if err != nil {
					lastErr = err
					continue
				}
				_ = resp.Body.Close()

				if resp.StatusCode >= 200 && resp.StatusCode < 400 {
					fmt.Printf("[OK] %s is UP at %s (status %d)\n", svc.Name, ip, resp.StatusCode)
					up = true
					break
				}
				lastErr = fmt.Errorf("status %d", resp.StatusCode)
			}

			prev, seen := wasUp[svc.Name]

			if !up {
				fmt.Printf("[FAIL] %s is DOWN: %v\n", svc.Name, lastErr)
				if !seen || prev {
					sendDiscordAlert(client, fmt.Sprintf("🔴 **%s** is DOWN: %v", svc.Name, lastErr))
				}
			} else if seen && !prev {
				sendDiscordAlert(client, fmt.Sprintf("🟢 **%s** has recovered", svc.Name))
			}

			wasUp[svc.Name] = up
		}
	}
}

func sendDiscordAlert(client *http.Client, message string) {
	webhookURL := os.Getenv("DISCORD_WEBHOOK")
	if webhookURL == "" {
		fmt.Println("DISCORD_WEBHOOK not set, skipping alert")
		return
	}

	payload := map[string]string{"content": message}
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling discord payload: %v\n", err)
		return
	}

	resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Error sending discord alert: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("Discord webhook returned status %d\n", resp.StatusCode)
	}
}
