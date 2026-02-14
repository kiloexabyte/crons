# crons

A custom cron scheduler for Home Assistant light control and service health monitoring.

## Features

- **Light Scheduling**: Automatically dim lights at night (8 PM → 10%) and brighten in the morning (5 AM → 100%)
- **Heater Control**: Turn heater on at 5 AM and off at 6 AM
- **Health Checks**: Monitor NAS services (Jellyfin, Sonarr, Radarr) every 5 minutes with Discord alerts on failure/recovery

## Setup

1. Copy `.env.example` to `.env` and configure:
   ```
   HA_TOKEN=your_home_assistant_token
   DISCORD_WEBHOOK=your_discord_webhook_url
   ```

2. Build and deploy to Raspberry Pi:
   ```bash
   op deploy
   ```

## Op Commands

| Command | Description |
|---------|-------------|
| `op build` | Build for Raspberry Pi (linux/arm64) |
| `op deploy` | Build and deploy to Pi |
| `op lint` | Run linters |
| `BRIGHTNESS=50 op testlights` | Test lights at specified brightness (0-100) |
| `STATE=on op testheater` | Test heater (on/off) |

## Requirements

- Go 1.24+
- [op](https://github.com/lesiw/op) CLI tool
