package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

var (
	Version   = "dev"
	CommitID  = "unknown"
	BuildDate = "unknown"
)

type Config struct {
	ServerURL string `json:"server_url"`
	Token     string `json:"token"`
	Share     bool   `json:"share"`
}

func main() {
	binaryName := filepath.Base(os.Args[0])

	var (
		showHelp    bool
		showVersion bool
		serverURL   string
		token       string
		share       string
		graph       string
	)

	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.StringVar(&serverURL, "server", "", "Server URL")
	flag.StringVar(&token, "token", "", "API token")
	flag.StringVar(&share, "share", "true", "Enable share link (true/false)")
	flag.StringVar(&graph, "graph", "", "Show graph for date range (YYYY-MM-DD:YYYY-MM-DD)")

	flag.Usage = func() {
		fmt.Printf(`%s - casspeed CLI Client

Usage: %s [options]

Options:
  --help              Show this help
  --version           Show version
  --server URL        Server URL (default: http://localhost:64580)
  --token TOKEN       API token for authenticated tests
  --share BOOL        Enable share link (default: true)
  --graph DATERANGE   Show historical graph (format: 2025-01-01:2025-01-31)

Examples:
  %s
  %s --server https://speed.example.com
  %s --token abc123 --share false
  %s --graph 2025-12-01:2025-12-31

`, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName)
	}

	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("%s version %s\n", binaryName, Version)
		fmt.Printf("Commit: %s\n", CommitID)
		fmt.Printf("Built: %s\n", BuildDate)
		fmt.Printf("Go: %s\n", runtime.Version())
		fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if serverURL == "" {
		serverURL = "http://localhost:64580"
	}

	if graph != "" {
		fmt.Println("📊 Historical graph not yet implemented")
		os.Exit(0)
	}

	fmt.Println("╭─────────────────────────────────────────────────╮")
	fmt.Println("│  ⚡ casspeed - Speed Test                       │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Printf("│  Server: %-39s│\n", serverURL)
	fmt.Println("╰─────────────────────────────────────────────────╯")
	fmt.Println()

	if err := runTest(serverURL, token, share == "true"); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func runTest(serverURL, token string, enableShare bool) error {
	u, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("invalid server URL: %w", err)
	}

	wsScheme := "ws"
	if u.Scheme == "https" {
		wsScheme = "wss"
	}

	wsURL := fmt.Sprintf("%s://%s/api/v1/speedtest/ws", wsScheme, u.Host)
	if !enableShare {
		wsURL += "?share=false"
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("connecting to server: %w", err)
	}
	defer conn.Close()

	var downloadSpeed, uploadSpeed, pingMs float64
	lastStage := ""

	for {
		var update map[string]interface{}
		if err := conn.ReadJSON(&update); err != nil {
			break
		}

		stage, _ := update["stage"].(string)
		progress, _ := update["progress"].(float64)
		speed, _ := update["speed"].(float64)
		message, _ := update["message"].(string)

		if stage != lastStage {
			lastStage = stage
			switch stage {
			case "ping":
				fmt.Println("🏓 Testing ping...")
			case "download":
				fmt.Println("⬇️  Testing download...")
			case "upload":
				fmt.Println("⬆️  Testing upload...")
			}
		}

		if progress > 0 {
			bar := makeProgressBar(int(progress * 50))
			fmt.Printf("\r%s %.0f%%  %s", bar, progress*100, message)
		}

		if stage == "ping" && progress >= 1.0 {
			pingMs = speed
			fmt.Println()
		} else if stage == "download" && progress >= 1.0 {
			downloadSpeed = speed
			fmt.Println()
		} else if stage == "upload" && progress >= 1.0 {
			uploadSpeed = speed
			fmt.Println()
		}

		if stage == "complete" {
			fmt.Println()
			break
		}
	}

	fmt.Println()
	fmt.Println("╭─────────────────────────────────────────────────╮")
	fmt.Println("│  ✅ Results                                     │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Printf("│  Download: %-37.1f Mbps │\n", downloadSpeed)
	fmt.Printf("│  Upload:   %-37.1f Mbps │\n", uploadSpeed)
	fmt.Printf("│  Ping:     %-37.1f ms   │\n", pingMs)
	fmt.Println("╰─────────────────────────────────────────────────╯")

	return nil
}

func makeProgressBar(width int) string {
	bar := "["
	for i := 0; i < 50; i++ {
		if i < width {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"
	return bar
}
