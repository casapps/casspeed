package swagger

import (
	"fmt"
	"net/http"
)

// Handler serves the Swagger UI
func Handler(w http.ResponseWriter, r *http.Request) {
	theme := DetectTheme(r)
	
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>casspeed API Documentation</title>
	<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
	<style>
		%s
	</style>
</head>
<body>
	<div id="swagger-ui"></div>
	<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
	<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
	<script>
		window.onload = function() {
			window.ui = SwaggerUIBundle({
				url: '/openapi.json',
				dom_id: '#swagger-ui',
				deepLinking: true,
				presets: [
					SwaggerUIBundle.presets.apis,
					SwaggerUIStandalonePreset
				],
				layout: "StandaloneLayout"
			});
		};
	</script>
</body>
</html>
`, GetThemeCSS(theme))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// SpecHandler serves the OpenAPI specification JSON
func SpecHandler(w http.ResponseWriter, r *http.Request) {
	spec := `{
	"openapi": "3.0.0",
	"info": {
		"title": "casspeed API",
		"description": "Self-hosted speed testing server API",
		"version": "1.0.0"
	},
	"servers": [
		{
			"url": "/api/v1",
			"description": "API v1"
		}
	],
	"paths": {
		"/speedtest/start": {
			"post": {
				"summary": "Start speed test",
				"description": "Initialize a new speed test",
				"responses": {
					"200": {
						"description": "Test started",
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"test_id": {
											"type": "string"
										},
										"status": {
											"type": "string"
										}
									}
								}
							}
						}
					}
				}
			}
		},
		"/speedtest/ws": {
			"get": {
				"summary": "Speed test WebSocket",
				"description": "WebSocket endpoint for real-time test progress"
			}
		},
		"/speedtest/download": {
			"get": {
				"summary": "Download test endpoint",
				"description": "Generates random data for download speed testing"
			}
		},
		"/speedtest/upload": {
			"post": {
				"summary": "Upload test endpoint",
				"description": "Consumes uploaded data for upload speed testing"
			}
		},
		"/speedtest/result/{id}": {
			"get": {
				"summary": "Get test result",
				"description": "Retrieve specific test result by ID",
				"parameters": [
					{
						"name": "id",
						"in": "path",
						"required": true,
						"schema": {
							"type": "string"
						}
					}
				]
			}
		},
		"/speedtest/history": {
			"get": {
				"summary": "Get test history",
				"description": "List historical test results"
			}
		}
	}
}
`
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, spec)
}
