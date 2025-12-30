# Installation Guide

## Docker (Recommended)

### Using Docker Compose

1. Clone the repository:
```bash
git clone https://github.com/casapps/casspeed.git
cd casspeed
```

2. Start the server:
```bash
docker-compose up -d
```

3. Access the web UI:
```
http://localhost:64580
```

### Using Docker CLI

```bash
docker run -d \
  -p 64580:64580 \
  -v casspeed-data:/var/lib/casapps/casspeed \
  ghcr.io/casapps/casspeed:latest
```

## Binary Installation

### Download Pre-built Binary

Download from [GitHub Releases](https://github.com/casapps/casspeed/releases):

```bash
# Linux AMD64
wget https://github.com/casapps/casspeed/releases/latest/download/casspeed-linux-amd64
chmod +x casspeed-linux-amd64
sudo mv casspeed-linux-amd64 /usr/local/bin/casspeed

# macOS ARM64
wget https://github.com/casapps/casspeed/releases/latest/download/casspeed-darwin-arm64
chmod +x casspeed-darwin-arm64
sudo mv casspeed-darwin-arm64 /usr/local/bin/casspeed
```

### Build from Source

Requirements:
- Go 1.23 or higher
- Make

```bash
git clone https://github.com/casapps/casspeed.git
cd casspeed
make build
sudo cp binaries/casspeed /usr/local/bin/
sudo cp binaries/casspeed-cli /usr/local/bin/
```

## CLI Client Installation

Download the CLI client:

```bash
# Linux AMD64
wget https://github.com/casapps/casspeed/releases/latest/download/casspeed-cli-linux-amd64
chmod +x casspeed-cli-linux-amd64
sudo mv casspeed-cli-linux-amd64 /usr/local/bin/casspeed-cli

# macOS ARM64
wget https://github.com/casapps/casspeed/releases/latest/download/casspeed-cli-darwin-arm64
chmod +x casspeed-cli-darwin-arm64
sudo mv casspeed-cli-darwin-arm64 /usr/local/bin/casspeed-cli
```

## Systemd Service (Linux)

```bash
sudo casspeed --service install
sudo systemctl enable casspeed
sudo systemctl start casspeed
```

## Configuration

On first run, casspeed creates a default configuration:

- Linux: `~/.config/casapps/casspeed/server.yml`
- macOS: `~/.config/casapps/casspeed/server.yml`
- Windows: `%APPDATA%\casapps\casspeed\server.yml`

Edit the configuration file and restart the service.
