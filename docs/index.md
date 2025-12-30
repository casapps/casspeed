# casspeed

**Self-hosted Speed Testing Server**

casspeed is a free, open-source alternative to speedtest.net that you can host yourself. It provides complete feature parity with commercial speed testing services while maintaining your privacy and giving you full control over your data.

## Features

- **Web-based speedtest interface** with real-time graphs
- **CLI client** for server-side speed testing  
- **Historical test results** tracking and visualization
- **No advertisements** or tracking
- **Self-hosted** for complete privacy and control
- **Multi-threaded** download/upload tests
- **Latency and jitter** measurements
- **Multiple concurrent tests** support
- **Share results** with generated images (PNG/SVG)
- **Rate limiting** and concurrent test management
- **SQLite database** for result storage

## Quick Start

### Docker (Recommended)

```bash
docker run -d \
  -p 64580:80 \
  -v ./config:/config \
  -v ./data:/data \
  ghcr.io/casapps/casspeed:latest
```

### Binary

```bash
# Download latest release
curl -LO https://github.com/casapps/casspeed/releases/latest/download/casspeed-linux-amd64

# Make executable
chmod +x casspeed-linux-amd64

# Run
./casspeed-linux-amd64 --port 8080
```

## Documentation

- [Installation Guide](installation.md) - Detailed installation instructions
- [API Reference](api.md) - REST API documentation
- [Configuration](configuration.md) - Configuration options

## Support

- GitHub: [casapps/casspeed](https://github.com/casapps/casspeed)
- Issues: [GitHub Issues](https://github.com/casapps/casspeed/issues)

## License

MIT License - see LICENSE.md for details
