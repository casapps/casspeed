# API Documentation

## Base URL

```
http://localhost:64580/api/v1
```

## Endpoints

### Speed Test

#### Start Test
```
POST /api/v1/speedtest/start
```

Response:
```json
{
  "test_id": "abc123",
  "status": "started"
}
```

#### WebSocket Progress
```
GET /api/v1/speedtest/ws
```

Real-time progress updates via WebSocket:
```json
{
  "stage": "download",
  "progress": 0.5,
  "speed": 123.4,
  "message": "123.4 Mbps"
}
```

#### Download Test
```
GET /api/v1/speedtest/download
```

Returns random data for download speed testing.

#### Upload Test
```
POST /api/v1/speedtest/upload
```

Accepts data upload for speed testing.

#### Get Result
```
GET /api/v1/speedtest/result/{id}
```

Response:
```json
{
  "id": "abc123",
  "download_mbps": 123.4,
  "upload_mbps": 56.7,
  "ping_ms": 12.3,
  "timestamp": "2025-12-28T00:00:00Z"
}
```

#### Get History
```
GET /api/v1/speedtest/history?user_id={id}
```

Returns array of test results.

### Share

#### View Share
```
GET /share/{code}
GET /s/{code}
```

HTML page with test results.

#### Share Images
```
GET /share/{code}.png
GET /s/{code}.png
```

PNG image (1200x630) for social media.

```
GET /share/{code}.svg
GET /s/{code}.svg
```

SVG image (scalable vector).

## Authentication

Use API tokens for authenticated requests:

```
Authorization: Bearer YOUR_TOKEN
```
