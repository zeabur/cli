# Zeabur CLI - Development Notes

## Build & Test
- Build: `go build ./...`
- Run: `go run ./cmd/main.go <command>`
- Test: `go test ./...`

## Project Structure
- `cmd/main.go` — entry point
- `internal/cmd/<command>/` — each CLI command in its own package
- `internal/cmdutil/` — shared command utilities (Factory, auth checks, spinner config)
- `pkg/api/` — GraphQL API client
- `pkg/model/` — data models (GraphQL struct tags)
- `internal/cmd/root/root.go` — root command, registers all subcommands

## Important: Keep `help --all` in sync
When adding or modifying CLI commands, flags, or subcommands, the output of `zeabur help --all` automatically reflects changes (it walks the Cobra command tree at runtime). No manual update is needed for the help output itself.

However, when adding a **new subcommand**, you must:
1. Create the command package under `internal/cmd/<parent>/<new>/`
2. Register it in the parent command file (e.g., `internal/cmd/template/template.go`)

## Conventions
- Each subcommand lives in its own package: `internal/cmd/<parent>/<sub>/<sub>.go`
- Commands support both interactive and non-interactive modes; if a flag is provided, skip the interactive prompt
- Use `cmdutil.SpinnerCharSet`, `cmdutil.SpinnerInterval`, `cmdutil.SpinnerColor` for spinners
- Models in `pkg/model/` use `graphql:"fieldName"` struct tags — only add fields that exist in the backend GraphQL schema
- Backend GraphQL schema lives in `../backend/internal/gateway/graphql/`
