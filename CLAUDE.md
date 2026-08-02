# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

shinpuru is a multi-purpose Discord bot written in Go, using `bwmarrin/discordgo` for the Discord API/gateway and `zekrotja/ken` as the slash command framework. It ships with a REST API (fiber) and a React SPA web frontend served by the same binary. The project uses [Task](https://taskfile.dev) (`Taskfile.yml`) as its command runner — most day-to-day commands below are `task <name>`.

## Common Commands

### Backend (Go)
- `task build-be` — build the backend binary to `bin/shinpuru`.
- `task run` — build and run the backend using `config/private.config.yml` (requires `task init-dev` first to create that file from `config/my.private.config.yml`).
- `task test` — run all backend unit tests: `go test -race -v -cover ./...`.
- Run a single package's tests: `go test -race -v -cover ./pkg/permissions/...`. Run a single test: `go test -race -v -run TestName ./pkg/permissions/...`.
- `task lint-be` — run `staticcheck` (must be installed: `go install honnef.co/go/tools/cmd/staticcheck@latest`).
- `task refresh-interfaces` — regenerate the `ISession` discordgo interface via `schnittstelle` (needed when the discordgo dependency version changes).
- `task refresh-mocks` — regenerate mocks in `mocks/` via `mockery` for all service interfaces (Database, ConfigProvider, KarmaProvider, IKen, etc). Run this after changing any interface listed in the Taskfile's `refresh-mocks` task.
- `task apidocs` — regenerate REST API docs from controller annotations (`swag` + `swagger-markdown` required).
- `task build-setup-tool` / `task build-cmdman-tool` — build the `cmd/setup` and `cmd/cmdman` helper binaries.

### Frontend (web/, React + Vite)
- `task run-fe-new` (or `cd web && yarn start`) — run the Vite dev server for the React app.
- `task build-fe` (or `cd web && yarn run build --base=/`) — production build; runs `tsc` then `vite build`.
- `task lint-fe` (or `cd web && yarn lint`) — ESLint.
- `task embed-fe` — build the frontend and copy it into `internal/util/embedded/webdist` so it gets embedded into the Go binary.
- `task release` — full release build: embeds frontend, builds frontend + backend + setup tool into `./release`.

Note: there is legacy Angular frontend tooling referenced in the Taskfile (`run-fe`, `deps-fe`) — the active frontend is the React app under `web/` (Vite-based, `run-fe-new`/`start` script).

## Architecture

### Dependency injection
Service wiring happens in `cmd/shinpuru/main.go` using `sarulabs/di`. Every service is registered once against a string identifier defined in `internal/util/static/di.go` (e.g. `static.DiDatabase`, `static.DiConfig`) and built lazily as an `App`-scoped singleton. Any component that needs a service takes a `di.Container` in its constructor and pulls dependencies out by key, e.g.:
```go
db := container.Get(static.DiDatabase).(database.Database)
```
Service dependencies must form a tree, not a cycle — a service built while resolving service A must not itself depend on A.

### Command handling (slash commands)
Slash commands live in `internal/slashcommands/`, one file per command, implementing `ken.SlashCommand` (`Name`, `Description`, `Version`, `Type`, `Options`, `Run`) and optionally `permissions.PermCommand` for permission-gated commands. New commands must be registered in `internal/inits/commandhandler.go` (`InitCommandHandler()`). Message-based and user-context commands live in `internal/messagecommands/` and `internal/usercommands/` respectively. Command-level middleware (logging, stats, disabling) is in `internal/middleware/`.

### Discord event listeners
Event handlers live in `internal/listeners/` as structs exposing handler methods, constructed with the `di.Container`. They're registered in `internal/inits/botsession.go` (`InitDiscordBotSession()`) via `session.AddHandler(listeners.NewXxx(container).Handler)`.

### State management
Discord state/caching (users, guilds, channels, sharding) is handled by `zekrotja/dgrs`, a Redis-backed state manager (not discordgo's built-in in-memory cache), registered under `static.DiState`.

### Database
`internal/services/database/database.go` defines the `Database` interface used everywhere data access is needed. Drivers: `internal/services/database/mysql` (primary, MySQL/MariaDB) and `internal/services/database/redis` (a caching middleware that wraps another `Database` driver — read-through/write-through: checks Redis first, falls back to the wrapped DB, keeps both hot). When adding functionality, extend the `Database` interface, the MySQL driver, and (if cacheable) the Redis middleware together. Schema changes to existing tables go through `internal/services/database/mysql/migrations.go`; brand-new tables are added directly to `setup()` in `mysql.go`.

### REST API / web server
`internal/services/webserver` hosts a fiber-based HTTP server. Structure:
- `v1/router.go` — versioned router (`/api/v1`).
- `v1/controllers/` — one controller per logical resource (guilds, backups, members, ...); each controller has a `Setup` method where its routes and required middleware are registered, and pulls its own dependencies from the `di.Container`.
- `v1/models/` — request/response DTOs plus transform functions from internal/discordgo types (e.g. `discordgo.Guild` → `models.Guild`).
- `middleware/` — endpoint- and controller-level middleware (auth token checks, permission checks).
- `auth/` — token/OAuth handling.
Three middleware tiers exist: global (rate limiting, CORS, filesystem), controller-level (mainly auth), and endpoint-level (mainly permission checks).

### Scheduler
`internal/services/scheduler` wraps `robfig/cron` (`CronLifeCycleWrapper`) for periodic jobs (expired votes, expired token cleanup, guild backups). Jobs are registered in `internal/inits/scheduler.go` using crontab syntax.

### Storage
`internal/services/storage` defines the `Storage` interface for object storage (avatars, images, backup files), with `file` (local disk) and `minio` (S3-compatible) drivers.

### Permissions
Guild-level command permissions are managed by `internal/services/permissions` (`Provider` interface) together with `pkg/permissions` (permission array encoding). Commands implement `permissions.PermCommand` to declare their default permission and domain.

### Utility packages (`pkg/`)
Framework-agnostic, independently reusable, well-tested packages — check here before writing new low-level helpers. Notable ones: `fetch` (fuzzy lookup of users/members/roles by ID/name/mention), `discordutil` (cache-first Discord object retrieval, permission checks), `embedbuilder` (embed builder pattern), `acceptmsg` (button-confirmation embeds), `permissions` (permission bit arrays). Shinpuru-specific shared helpers that need internal types live in `internal/util/` instead.

### Mocks
Interface mocks for tests live in `mocks/` and are generated by `mockery` via `task refresh-mocks` — do not hand-edit them; update the source interface and regenerate instead.

### Web frontend (`web/`)
React + TypeScript SPA built with Vite, styled with `styled-components`, global state via `zustand`, localized with `i18next` (locale JSON files in `web/public/locales`, one namespace per route — start there for text changes). API access goes through the custom client in `web/src/lib/shinpuru-ts` via the `useApi` hook (`web/src/hooks/useApi.ts`); other reusable hooks also live in `web/src/hooks`. The built app is embedded into the Go binary at `internal/util/embedded/webdist` (via `task embed-fe`) and served directly by the backend's web server.

### Configuration
shinpuru is configured via a YAML/JSON file passed with `-c`. See `config/config.example.yml` for all documented keys. For local dev, copy `config/my.private.config.yml` to `config/private.config.yml` (via `task init-dev`) and fill in credentials — this file is gitignored.
