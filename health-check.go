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

var nasServices = []struct {
	Name string
	Port string
}{
	{"Jellyfin", "8096"},
	{"Sonarr", "8989"},
	{"Radarr", "8310"},
}

var directServices = []struct {
	Name string
	URL  string
}{
	{"Home Assistant", "http://homeassistant:8123/"},
}

func checkServices() Action {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Track previous state for each service
	wasUp := make(map[string]bool)

	return func() {
		// Check NAS services (with IP fallback)
		for _, svc := range nasServices {
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

			checkResult(client, wasUp, svc.Name, up, lastErr)
		}

		// Check direct URL services
		for _, svc := range directServices {
			up := false
			var lastErr error

			resp, err := client.Get(svc.URL)
			if err != nil {
				lastErr = err
			} else {
				_ = resp.Body.Close()
				if resp.StatusCode >= 200 && resp.StatusCode < 400 {
					fmt.Printf("[OK] %s is UP (status %d)\n", svc.Name, resp.StatusCode)
					up = true
				} else {
					lastErr = fmt.Errorf("status %d", resp.StatusCode)
				}
			}

			checkResult(client, wasUp, svc.Name, up, lastErr)
		}
	}
}

func checkResult(client *http.Client, wasUp map[string]bool, name string, up bool, lastErr error) {
	prev, seen := wasUp[name]

	if !up {
		fmt.Printf("[FAIL] %s is DOWN: %v\n", name, lastErr)
		if !seen || prev {
			sendDiscordAlert(client, fmt.Sprintf("🔴 **%s** is DOWN: %v", name, lastErr))
		}
	} else if seen && !prev {
		sendDiscordAlert(client, fmt.Sprintf("🟢 **%s** has recovered", name))
	}

	wasUp[name] = up
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
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		fmt.Printf("Discord webhook returned status %d\n", resp.StatusCode)
	}
}
