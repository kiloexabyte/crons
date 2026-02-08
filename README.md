# crons

A custom cron scheduler for Home Assistant light control and service health monitoring.

## Features

- **Light Scheduling**: Automatically dim lights at night and brighten in the morning
- **Health Checks**: Monitor NAS services (Jellyfin, Sonarr, Radarr) with Discord alerts on failure/recovery

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
| `op testlights` | Test lights (`BRIGHTNESS=50 op testlights`) |
| `op testdiscord` | Test Discord webhook |

## Requirements

- Go 1.24+
- [op](https://github.com/lesiw/op) CLI tool
