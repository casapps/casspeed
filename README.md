## 👋 Welcome to casspeed 🚀

A free self-hosted alternative to speedtest.net with all the features but free and opensource with no ads, tracking, and no feature gating.

## Features

- 🚀 Multi-threaded download/upload tests
- 📊 Real-time WebSocket progress updates
- 🔗 Shareable test results with PNG/SVG export
- 👥 Multi-user support with device tracking
- 🔐 API token authentication
- 📱 Responsive dark theme web UI
- 💻 CLI client with real-time display and graphs
- 🐳 Docker and multi-platform support

## Quick Start

### Docker

```bash
docker-compose up -d
open http://localhost:64580
```

### Build from Source

```bash
make build
./bin/casspeed
```

### CLI Client

```bash
casspeed-cli --token YOUR_TOKEN
casspeed-cli --graph 2025-12-01:2025-12-31
```

## API Endpoints

- `GET /` - Web UI
- `POST /api/v1/speedtest/start` - Start test
- `GET /api/v1/speedtest/ws` - WebSocket progress
- `GET /share/{code}` - View shared result
- `GET /share/{code}.png` - PNG image
- `GET /share/{code}.svg` - SVG image

## Configuration

Configuration file: `server.yml` (auto-created on first run)

```yaml
server:
  port: 64580
  mode: production

test:
  max_concurrent: 3
  default_duration: 10
  max_threads: 16
```

## Author

🤖 casjay: [Github](https://github.com/casjay) 🤖

