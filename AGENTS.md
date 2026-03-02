# Repository Guidelines

## Project Structure & Module Organization
`NodeConverter` is a Go service that converts proxy subscription formats.

- `main.go`: application entrypoint and HTTP server bootstrap.
- `core/`: conversion logic and protocol models (`vless`, `trojan`, `ss`, `tuic`, `hysteria`, etc.).
- `handler/`: HTTP handlers, fetch/filter flow, and request-level behavior.
- `lib/`: shared utilities (network helpers, YAML emoji handling).
- `config.yaml` and `clash-tpl.yaml`: runtime and template configuration.
- `.github/workflows/go.yml`: release workflow (GoReleaser on tag push).

## Build, Test, and Development Commands
- `go run .`  
  Run the service locally (default config from repository root).
- `go build ./...`  
  Compile all packages and catch build-time issues.
- `go test ./...`  
  Run all unit tests across `core/` and `handler/`.
- `docker compose up --build`  
  Start the containerized service for integration-style local checks.

## Coding Style & Naming Conventions
- Follow standard Go formatting: run `gofmt -w .` before submitting.
- Keep package names short, lowercase, and domain-oriented (`core`, `handler`, `lib`).
- Exported identifiers: `PascalCase`; internal helpers: `camelCase`.
- Prefer focused files by protocol/feature (for example `core/vless.go`, `core/vless_test.go`).

## Testing Guidelines
- Use Go’s built-in `testing` package.
- Place tests next to implementation files with `_test.go` suffix.
- Use explicit test names like `TestParseVLESS` and table-driven tests where practical.
- Ensure `go test ./...` passes before opening a pull request.

## Commit & Pull Request Guidelines
- Existing history mixes Conventional Commits and concise fix-style messages. Prefer:
  - `feat(scope): add ...`
  - `fix(scope): correct ...`
  - `refactor(scope): simplify ...`
- Keep each commit focused on one change area.
- PRs should include:
  - Clear summary of behavior changes.
  - Linked issue (if applicable).
  - Test evidence (`go test ./...` output or equivalent).
  - Example request/response when API behavior changes (for `/sub` parameters).
