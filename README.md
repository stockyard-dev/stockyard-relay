# Stockyard Relay

**Webhook sender — define events, external services register callbacks, fire and deliver with retry**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9090:9090 -v relay_data:/data ghcr.io/stockyard-dev/stockyard-relay
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9090` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9090` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `RELAY_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 5 events, 10 subscriptions | Unlimited events and subscriptions |
| Price | Free | $2.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Developer Tools

## License

Apache 2.0
