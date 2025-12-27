# TODO List for casspeed Project

## Phase 1: Foundation ✓
- [x] Read entire TEMPLATE.md
- [x] Create AI.md (1.1MB, 29,874 lines)
- [x] Fill PART 36 with speedtest business logic
- [x] Create directory structure
- [x] Create go.mod and release.txt

## Phase 2: Infrastructure (Makefile, Docker, CI/CD) ✓
- [x] Create Makefile with build/release/docker/test/dev targets
- [x] Create Dockerfile (multi-stage, builder + runtime)
- [x] Create entrypoint.sh with signal handling
- [x] Create docker-compose.yml (production)
- [x] Create docker-compose.dev.yml (development)
- [x] Create docker-compose.test.yml (testing)
- [x] Create .dockerignore
- [x] Create GitHub Actions release.yml workflow
- [x] Create GitHub Actions docker.yml workflow

## Phase 3: Core Packages
- [x] Implement config/bool.go
- [ ] Read & implement config/config.go (PART 5)
- [ ] Read & implement mode/mode.go (PART 6)
- [ ] Read & implement paths/paths.go (PART 4)
- [ ] Read & implement ssl/ssl.go (PART 21)
- [ ] Read & implement scheduler/scheduler.go (PART 27)
- [ ] Read & implement service/service.go (PART 10)

## Phase 4: Server Binary & Main
- [ ] Read PART 7 completely
- [ ] Implement src/main.go with all CLI flags
- [ ] Implement --help, --version handlers
- [ ] Implement --service commands (start/stop/restart/etc.)
- [ ] Implement --maintenance commands
- [ ] Implement --update command
- [ ] Implement --daemon mode
- [ ] Implement server/server.go HTTP setup

## Phase 5: Database & Models
- [ ] Read PART 23, PART 24
- [ ] Implement store/store.go interface
- [ ] Implement store/sqlite.go (modernc.org/sqlite)
- [ ] Create migrations system
- [ ] Create users table schema
- [ ] Create devices table schema  
- [ ] Create speed_tests table schema
- [ ] Create sessions table schema
- [ ] Create api_tokens table schema
- [ ] Implement all model/ structs

## Phase 6: User Management & Auth
- [ ] Implement service/user.go
- [ ] Implement service/device.go
- [ ] Implement service/auth.go (Argon2id)
- [ ] Implement service/token.go
- [ ] Implement handler/auth.go
- [ ] Implement handler/user.go
- [ ] Implement TOTP/Passkeys support
- [ ] Implement OIDC/LDAP support

## Phase 7: Speedtest Engine & API
- [ ] Implement service/speedtest.go
- [ ] Multi-threaded download test logic
- [ ] Multi-threaded upload test logic
- [ ] Ping/latency test logic
- [ ] WebSocket progress tracking
- [ ] Implement handler/speedtest.go
- [ ] POST /api/v1/speedtest/start
- [ ] GET /api/v1/speedtest/download
- [ ] POST /api/v1/speedtest/upload
- [ ] GET /api/v1/speedtest/ping
- [ ] WebSocket /api/v1/speedtest/status/:id
- [ ] GET /api/v1/speedtest/result/:id
- [ ] GET /api/v1/speedtest/results
- [ ] GET /api/v1/speedtest/stats
- [ ] DELETE /api/v1/speedtest/result/:id
- [ ] GET /api/v1/speedtest/export
- [ ] GET /api/v1/speedtest/server
- [ ] Rate limiting middleware

## Phase 8: Admin Panel
- [ ] Read PART 19
- [ ] Implement admin/admin.go
- [ ] Admin authentication
- [ ] User management UI
- [ ] Device management UI
- [ ] Server settings UI
- [ ] Test results overview

## Phase 9: Web Frontend
- [ ] Read PART 17
- [ ] Create HTML templates
- [ ] Embed static assets (CSS, JS)
- [ ] Light/dark theme CSS
- [ ] Speedtest UI with real-time display
- [ ] Results history page with graphs
- [ ] User registration/login pages
- [ ] Device management page

## Phase 10: API Documentation
- [ ] Read PART 20 swagger section
- [ ] Implement swagger/swagger.go
- [ ] Swagger UI handler
- [ ] API annotations
- [ ] Read PART 20 graphql section
- [ ] Implement graphql/graphql.go
- [ ] GraphQL schema
- [ ] GraphQL resolvers
- [ ] GraphiQL UI

## Phase 11: CLI Client (REQUIRED)
- [ ] Read PART 34 completely
- [ ] Implement src/client/main.go
- [ ] CLI command: test
- [ ] CLI command: history
- [ ] CLI command: export
- [ ] CLI command: login
- [ ] CLI command: device
- [ ] Token authentication
- [ ] Device registration
- [ ] ASCII graph output
- [ ] CSV/JSON export
- [ ] Config file support

## Phase 12: Documentation
- [ ] Update README.md
- [ ] Read PART 33
- [ ] Create mkdocs.yml
- [ ] Create .readthedocs.yaml
- [ ] Create docs/index.md
- [ ] Create docs/installation.md
- [ ] Create docs/configuration.md
- [ ] Create docs/api.md
- [ ] Create docs/cli.md
- [ ] Create docs/stylesheets/dark.css

## Phase 13: Testing & Scripts
- [ ] Read PART 13
- [ ] Create tests/run_tests.sh
- [ ] Create tests/docker.sh
- [ ] Create tests/incus.sh
- [ ] Write unit tests for config/
- [ ] Write unit tests for service/
- [ ] Write integration tests

## Phase 14: Security & Features
- [ ] Implement backup/restore (PART 25)
- [ ] Implement email notifications (PART 26)
- [ ] Implement metrics (PART 29)
- [ ] Implement Tor hidden service (PART 30)
- [ ] Error handling (PART 31)

## Phase 15: Final Compliance
- [ ] Full compliance check against AI.md
- [ ] Test all 8 platforms build
- [ ] Test Docker deployment
- [ ] Test CLI client
- [ ] Test all API endpoints
- [ ] Test admin panel
- [ ] Test web UI
- [ ] Update documentation
- [ ] Empty TODO.AI.md
- [ ] Write completion COMMIT_MESS
