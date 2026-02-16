package homeassistant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const DefaultURL = "http://homeassistant:8123"

var AllLights = []string{
	"light.pc1",
	"light.pc2",
	"light.pc3",
	"light.kitchen1",
	"light.kitchen2",
	"light.kitchen3",
	"light.livingroom1",
	"light.livingroom2",
	"light.livingroom3",
	"light.bed",
}

type brightnessPayload struct {
	EntityID      string `json:"entity_id"`
	BrightnessPct int    `json:"brightness_pct,omitempty"`
}

type Client struct {
	URL   string
	Token string
}

func NewClient() *Client {
	return &Client{
		URL:   DefaultURL,
		Token: os.Getenv("HA_TOKEN"),
	}
}

func (c *Client) SetLightBrightness(entityID string, percent int) error {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	service := "turn_on"
	if percent == 0 {
		service = "turn_off"
	}

	url := fmt.Sprintf("%s/api/services/light/%s", c.URL, service)

	payload := brightnessPayload{EntityID: entityID}
	if percent > 0 {
		payload.BrightnessPct = percent
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error calling Home Assistant: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("home assistant returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) SetAllLightsBrightness(percent int) error {
	for _, entityID := range AllLights {
		if err := c.SetLightBrightness(entityID, percent); err != nil {
			fmt.Printf("Error setting %s: %v\n", entityID, err)
			continue
		}
		fmt.Printf("Light %s set to %d%%\n", entityID, percent)
	}
	return nil
}

type switchPayload struct {
	EntityID string `json:"entity_id"`
}

func (c *Client) SetSwitch(entityID string, on bool) error {
	service := "turn_off"
	if on {
		service = "turn_on"
	}

	url := fmt.Sprintf("%s/api/services/switch/%s", c.URL, service)

	body, _ := json.Marshal(switchPayload{EntityID: entityID})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error calling Home Assistant: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("home assistant returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) SetHeater(on bool) error {
	return c.SetSwitch("switch.heater", on)
}
