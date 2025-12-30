package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/casapps/casspeed/src/config"
	"github.com/casapps/casspeed/src/mode"
	"github.com/casapps/casspeed/src/paths"
	"github.com/casapps/casspeed/src/server"
)

// Version information (set by linker flags during build)
var (
	Version   = "dev"
	CommitID  = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Get binary name for help text
	binaryName := filepath.Base(os.Args[0])

	// Define flags
	var (
		showHelp    bool
		showVersion bool
		showStatus  bool
		daemonFlag  bool
		modeFlag    string
		debugFlag   string
		configDir   string
		dataDir     string
		logDir      string
		pidFile     string
		address     string
		portFlag    string
		serviceCmd  string
		maintCmd    string
		updateCmd   string
	)

	flag.BoolVar(&showHelp, "help", false, "Show help information")
	flag.BoolVar(&showHelp, "h", false, "Show help information (short)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information (short)")
	flag.BoolVar(&showStatus, "status", false, "Show status and health")
	flag.BoolVar(&daemonFlag, "daemon", false, "Daemonize (detach from terminal)")
	flag.StringVar(&modeFlag, "mode", "", "Application mode (production|development)")
	flag.StringVar(&debugFlag, "debug", "", "Enable debug mode")
	flag.StringVar(&configDir, "config", "", "Configuration directory")
	flag.StringVar(&dataDir, "data", "", "Data directory")
	flag.StringVar(&logDir, "log", "", "Log directory")
	flag.StringVar(&pidFile, "pid", "", "PID file path")
	flag.StringVar(&address, "address", "", "Listen address")
	flag.StringVar(&portFlag, "port", "", "Listen port")
	flag.StringVar(&serviceCmd, "service", "", "Service management (start|stop|restart|reload|install|uninstall|help)")
	flag.StringVar(&maintCmd, "maintenance", "", "Maintenance operations (backup|restore|update|mode|setup)")
	flag.StringVar(&updateCmd, "update", "", "Update operations (check|yes|branch stable|beta|daily)")

	flag.Usage = func() {
		showHelpText(binaryName)
	}

	flag.Parse()

	// Handle --help
	if showHelp {
		showHelpText(binaryName)
		os.Exit(0)
	}

	// Handle --version
	if showVersion {
		showVersionInfo(binaryName)
		os.Exit(0)
	}

	// Handle --status
	if showStatus {
		showStatusInfo(binaryName)
		os.Exit(0)
	}

	// Handle --service
	if serviceCmd != "" {
		handleService(binaryName, serviceCmd)
		os.Exit(0)
	}

	// Handle --maintenance
	if maintCmd != "" {
		handleMaintenance(binaryName, maintCmd, flag.Args())
		os.Exit(0)
	}

	// Handle --update
	if updateCmd != "" {
		handleUpdate(binaryName, updateCmd, flag.Args())
		os.Exit(0)
	}

	// Handle --daemon
	if daemonFlag {
		fmt.Fprintln(os.Stderr, "Error: --daemon not supported")
		fmt.Fprintln(os.Stderr, "Use systemd (Type=simple), Docker, or run in foreground")
		os.Exit(1)
	}

	// Detect application mode
	appMode, err := mode.Detect(modeFlag, debugFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Detect paths
	appPaths, err := paths.Detect(configDir, dataDir, logDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error detecting paths: %v\n", err)
		os.Exit(1)
	}

	// Set PID file if specified
	if pidFile != "" {
		appPaths.PID = pidFile
	}

	// Ensure all directories exist
	if err := appPaths.Ensure(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directories: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	configPath := filepath.Join(appPaths.Config, "server.yml")
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Override config with CLI flags
	if address != "" {
		cfg.Server.Address = address
	}
	if portFlag != "" {
		cfg.Server.Port = portFlag
	}
	if modeFlag != "" {
		cfg.Server.Mode = modeFlag
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Save default configuration if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := config.Save(cfg, configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not save default config: %v\n", err)
		}
	}

	// Print startup banner
	printBanner(appMode, cfg)

	// Determine listen address and port
	listenAddr := cfg.Server.Address
	if listenAddr == "" {
		listenAddr = "[::]"
	}

	listenPort := 64580
	if cfg.Server.Port != nil {
		switch p := cfg.Server.Port.(type) {
		case int:
			listenPort = p
		case string:
			if p != "" {
				fmt.Sscanf(p, "%d", &listenPort)
			}
		}
	}

	// Create and start server
	srv, err := server.New(cfg, appMode, appPaths.Data, Version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server initialization error: %v\n", err)
		os.Exit(1)
	}

	if err := srv.Start(listenAddr, listenPort); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func showHelpText(binaryName string) {
	fmt.Printf(`Usage: %s [options]

casspeed - Self-hosted Speed Testing Server

OPTIONS:
  -h, --help              Show this help message
  -v, --version           Show version information
  --status                Show status and health
  --mode MODE             Set application mode (production|development)
  --debug                 Enable debug mode (verbose logging, debug endpoints)
  --daemon                Daemonize (detach from terminal)
  --config DIR            Configuration directory
  --data DIR              Data directory  
  --log DIR               Log directory
  --pid FILE              PID file path
  --address ADDR          Listen address (default: [::])
  --port PORT             Listen port (default: random 64xxx)
  --service CMD           Service management (start|stop|restart|reload|install|uninstall|help)
  --maintenance CMD       Maintenance operations (backup|restore|update|mode|setup)
  --update CMD            Update operations (check|yes|branch stable|beta|daily)

EXAMPLES:
  %s
    Start server with defaults

  %s --port 8080
    Start on port 8080

  %s --mode production --port 443
    Start in production mode on port 443

  %s --config /etc/casspeed --data /var/lib/casspeed
    Start with custom directories

  %s --status
    Show server status

  %s --service start
    Start as service

  %s --maintenance backup
    Backup database

For more information, visit: https://github.com/casapps/casspeed
`, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName)
}

func showVersionInfo(binaryName string) {
	// Format: binaryname v1.2.3 (commit) built YYYY-MM-DD
	fmt.Printf("%s v%s (%s) built %s\n", binaryName, Version, CommitID, BuildDate)
}

func printBanner(appMode *mode.State, cfg *config.Config) {
	fmt.Println("╭─────────────────────────────────────────────────────────────╮")
	fmt.Printf("│  🚀 CASSPEED · 📦 v%s%s│\n", Version, pad(Version))
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	fmt.Printf("│  %s Running in mode: %s%s│\n", appMode.GetConsoleIcon(), appMode.String(), padMode(appMode.String()))
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	fmt.Println("│  Server initialization...                                   │")
	fmt.Println("╰─────────────────────────────────────────────────────────────╯")
}

func pad(version string) string {
	// Pad to align with banner width (60 chars minus prefix)
	needed := 60 - len("🚀 CASSPEED · 📦 v") - len(version)
	if needed < 0 {
		needed = 0
	}
	return fmt.Sprintf("%*s", needed, "")
}

func padMode(modeStr string) string {
	// Pad mode line
	needed := 60 - len("Running in mode: ") - len(modeStr) - 2
	if needed < 0 {
		needed = 0
	}
	return fmt.Sprintf("%*s", needed, "")
}

func showStatusInfo(binaryName string) {
	fmt.Printf("%s Status\n", binaryName)
	fmt.Println("─────────────────────────────────────")
	
	// Check if server is running via health endpoint
	resp, err := http.Get("http://localhost:64580/healthz")
	if err != nil {
		fmt.Println("Status: Stopped (not responding)")
		fmt.Println("Health: Unavailable")
		fmt.Println()
		fmt.Printf("Start the server with: %s\n", binaryName)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == 200 {
		fmt.Println("Status: Running")
		fmt.Println("Health: OK")
		fmt.Printf("Endpoint: http://localhost:64580\n")
	} else {
		fmt.Printf("Status: Running (unhealthy, HTTP %d)\n", resp.StatusCode)
		fmt.Println("Health: Error")
	}
}

func handleService(binaryName string, cmd string) {
	fmt.Printf("%s: Service Management\n", binaryName)
	fmt.Println("─────────────────────────────────────")
	
	switch cmd {
	case "start":
		fmt.Println("Service management: Use systemd, Docker, or run directly")
		fmt.Printf("Direct: %s --daemon\n", binaryName)
		fmt.Println("Systemd: systemctl start casspeed")
		fmt.Println("Docker: docker-compose up -d")
	case "stop":
		fmt.Println("Service management: Use systemd, Docker, or kill process")
		fmt.Println("Systemd: systemctl stop casspeed")
		fmt.Println("Docker: docker-compose down")
	case "restart":
		fmt.Println("Service management: Use systemd, Docker, or kill + restart")
		fmt.Println("Systemd: systemctl restart casspeed")
		fmt.Println("Docker: docker-compose restart")
	case "reload":
		fmt.Println("Configuration reload: Send SIGHUP to process")
		fmt.Println("Kill: pkill -HUP casspeed")
		fmt.Println("Systemd: systemctl reload casspeed")
	case "install", "--install":
		fmt.Println("Service installation:")
		fmt.Println("1. Copy binary to /usr/local/bin/casspeed")
		fmt.Println("2. Create systemd unit: /etc/systemd/system/casspeed.service")
		fmt.Println("3. Enable: systemctl enable casspeed")
		fmt.Println("Or use Docker Compose for containerized deployment")
	case "uninstall", "--uninstall":
		fmt.Println("Service removal:")
		fmt.Println("1. Stop service: systemctl stop casspeed")
		fmt.Println("2. Disable: systemctl disable casspeed")
		fmt.Println("3. Remove unit: rm /etc/systemd/system/casspeed.service")
		fmt.Println("4. Remove binary: rm /usr/local/bin/casspeed")
	case "help", "--help":
		fmt.Printf(`Service Management Commands:

  %s --service start       Start the service
  %s --service stop        Stop the service
  %s --service restart     Restart the service
  %s --service reload      Reload configuration
  %s --service install     Install service (systemd/launchd/etc)
  %s --service uninstall   Uninstall service
  %s --service help        Show this help

Note: Use systemd/launchd/Docker for production deployments.
Manual service management is for advanced users only.
`, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName, binaryName)
	default:
		fmt.Printf("Unknown service command: %s\n", cmd)
		fmt.Printf("Run '%s --service help' for available commands\n", binaryName)
		os.Exit(1)
	}
}

func handleMaintenance(binaryName string, cmd string, args []string) {
	fmt.Printf("%s: Maintenance Operations\n", binaryName)
	fmt.Println("─────────────────────────────────────")
	
	switch cmd {
	case "backup":
		fmt.Println("Database backup:")
		fmt.Println("Default DB: /var/lib/casapps/casspeed/db/speedtest.db")
		fmt.Println("Backup command: cp speedtest.db speedtest.db.backup")
		fmt.Println("Or use sqlite3 .backup command for consistency")
	case "restore":
		if len(args) > 0 {
			fmt.Printf("Database restore from: %s\n", args[0])
			fmt.Println("Restore command: cp backup.db speedtest.db")
			fmt.Println("Warning: Stop the server before restoring")
		} else {
			fmt.Println("Usage: --maintenance restore <backup-file>")
			os.Exit(1)
		}
	case "update":
		fmt.Println("Server update:")
		fmt.Println("1. Download latest binary from GitHub releases")
		fmt.Println("2. Stop server")
		fmt.Println("3. Replace binary")
		fmt.Println("4. Start server")
		fmt.Println()
		fmt.Println("Docker: docker-compose pull && docker-compose up -d")
	case "mode":
		if len(args) > 0 {
			mode := args[0]
			if mode == "production" || mode == "development" {
				fmt.Printf("Set mode in config: server.mode: %s\n", mode)
				fmt.Println("Or use environment: MODE=%s\n", mode)
				fmt.Printf("Or use CLI flag: %s --mode %s\n", binaryName, mode)
			} else {
				fmt.Println("Invalid mode. Use: production or development")
				os.Exit(1)
			}
		} else {
			fmt.Println("Usage: --maintenance mode <production|development>")
			os.Exit(1)
		}
	case "setup":
		fmt.Println("Setup guide:")
		fmt.Println("1. First run creates default config in /etc/casapps/casspeed/")
		fmt.Println("2. Edit server.yml for custom configuration")
		fmt.Println("3. Database auto-created on first run")
		fmt.Println("4. No additional setup required")
	default:
		fmt.Printf("Unknown maintenance command: %s\n", cmd)
		fmt.Printf("Available: backup, restore, update, mode, setup\n")
		os.Exit(1)
	}
}

func handleUpdate(binaryName string, cmd string, args []string) {
	fmt.Printf("%s: Update System\n", binaryName)
	fmt.Println("─────────────────────────────────────")
	
	switch cmd {
	case "check":
		fmt.Println("Checking for updates...")
		fmt.Printf("Current version: %s\n", Version)
		fmt.Println()
		fmt.Println("Check GitHub releases:")
		fmt.Println("https://github.com/casapps/casspeed/releases/latest")
	case "yes":
		fmt.Println("Update instructions:")
		fmt.Println("1. Download latest binary from GitHub releases")
		fmt.Println("2. Verify checksum (provided in release)")
		fmt.Println("3. Stop casspeed: systemctl stop casspeed")
		fmt.Println("4. Replace binary: cp casspeed-new /usr/local/bin/casspeed")
		fmt.Println("5. Start casspeed: systemctl start casspeed")
		fmt.Println()
		fmt.Println("Docker: docker-compose pull && docker-compose up -d")
	case "branch":
		if len(args) > 0 {
			branch := args[0]
			if branch == "stable" || branch == "beta" || branch == "daily" {
				fmt.Printf("Branch switching:\n")
				fmt.Println("Stable: Use tagged releases from GitHub")
				fmt.Println("Beta: Use pre-release tags")
				fmt.Println("Daily: Build from main branch")
				fmt.Println()
				fmt.Printf("Docker: ghcr.io/casapps/casspeed:%s\n", branch)
			} else {
				fmt.Printf("Invalid branch: %s (use: stable, beta, daily)\n", branch)
				os.Exit(1)
			}
		} else {
			fmt.Println("Usage: --update branch <stable|beta|daily>")
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown update command: %s\n", cmd)
		fmt.Printf("Available: check, yes, branch <stable|beta|daily>\n")
		os.Exit(1)
	}
}
