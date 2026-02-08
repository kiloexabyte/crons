package main

import (
	"fmt"
	"net/http"
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

			if !up {
				fmt.Printf("[FAIL] %s is DOWN: %v\n", svc.Name, lastErr)
				// TODO: send discord alert
			}
		}
	}
}
