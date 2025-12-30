# Configuration

casspeed uses a YAML configuration file located at `/etc/casapps/casspeed/server.yml` (or `~/.config/casapps/casspeed/server.yml` for non-root users).

## Configuration File Location

The configuration file location can be customized using the `--config` flag:

```bash
casspeed --config /path/to/config/dir
```

## Configuration Options

### Server Section

```yaml
server:
  # Listen port (default: random 64xxx)
  port: 64580
  
  # Listen address (default: [::] for IPv6 and IPv4)
  address: "[::]"
  
  # Fully qualified domain name
  fqdn: "speedtest.example.com"
  
  # Application mode (production or development)
  mode: production
  
  # Branding
  branding:
    title: "casspeed"
    tagline: "Self-hosted Speed Test"
    description: "Fast, private network speed testing"
  
  # SEO metadata
  seo:
    keywords:
      - speedtest
      - bandwidth
      - network
      - performance
```

### Test Section

```yaml
test:
  # Maximum concurrent tests per IP
  max_concurrent: 3
  
  # Minimum seconds between tests from same IP
  min_interval: 5
  
  # Default test duration in seconds
  default_duration: 10
  
  # Maximum threads for multi-threaded tests
  max_threads: 16
  
  # Days to keep test results (0 = unlimited)
  results_retention: 90
  
  # Data chunk size in bytes
  chunk_size: 1048576  # 1MB
  
  # Test timeout in seconds
  timeout: 60
```

### Web UI Section

```yaml
web:
  ui:
    # Theme: light, dark, or auto
    theme: dark
  
  # CORS configuration (* allows all origins)
  cors: "*"
```

### Rate Limiting

```yaml
server:
  rate_limit:
    # Enable rate limiting
    enabled: true
    
    # Maximum requests per window
    requests: 120
    
    # Window size in seconds
    window: 60
```

### Database Section

```yaml
server:
  database:
    # Driver: file (SQLite), postgres, mysql
    driver: file
    
    # For PostgreSQL/MySQL
    host: localhost
    port: 5432
    name: casspeed
    username: casspeed
    password: secret
    sslmode: disable
```

### SSL/TLS Section

```yaml
server:
  ssl:
    # Enable HTTPS
    enabled: false
    
    # Certificate and key paths
    cert: /etc/casspeed/ssl/cert.pem
    key: /etc/casspeed/ssl/key.pem
    
    # Minimum TLS version (TLS1.2 or TLS1.3)
    min_version: TLS1.2
    
    # Let's Encrypt configuration
    letsencrypt:
      enabled: false
      email: admin@example.com
      challenge: http-01  # http-01, tls-alpn-01, or dns-01
      staging: false
```

### Scheduler Section

```yaml
server:
  scheduler:
    enabled: true
    tasks:
      log_rotation:
        enabled: true
        schedule: "0 0 * * *"  # Daily at midnight
        max_age: "30d"
        max_size: "100MB"
      
      session_cleanup:
        enabled: true
        schedule: "@hourly"
      
      backup:
        enabled: true
        schedule: "0 2 * * *"  # Daily at 2 AM
        retention: 4  # Keep 4 backups
      
      ssl_renewal:
        enabled: true
        schedule: "0 3 * * *"  # Daily at 3 AM
        renew_before: "7d"  # Renew 7 days before expiry
      
      health_check:
        enabled: true
        schedule: "*/5 * * * *"  # Every 5 minutes
```

## Environment Variables

casspeed can be configured using environment variables in Docker:

```bash
docker run -d \
  -e MODE=production \
  -e PORT=80 \
  -e ADDRESS=0.0.0.0 \
  -e DEBUG=false \
  -e TZ=America/New_York \
  ghcr.io/casapps/casspeed:latest
```

## Default Values

If no configuration file exists, casspeed will create one with sane defaults:

- **Port:** Random in 64xxx range
- **Address:** `[::]` (all interfaces, IPv4 and IPv6)
- **Mode:** production
- **Max concurrent tests:** 3 per IP
- **Min interval:** 5 seconds
- **Theme:** dark
- **Rate limiting:** enabled (120 requests/minute)

## Validation

casspeed validates the configuration on startup and will exit with an error if any values are invalid.
