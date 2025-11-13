# Contributing to GreatestWorks

Thanks for your interest in contributing! This project implements a distributed MMO game server in Go using DDD and a layered microservice architecture. To keep quality high and changes consistent, please follow the guidelines below.

If you only read one extra document, read this one from the repo:
- .github/instructions/gw.instructions.md (authoritative engineering conventions, architecture, and dos/don'ts)

## Prerequisites
- Go 1.24+
- MongoDB 4.4+ (5.0+ recommended)
- Redis 6.0+ (7.0+ recommended)
- Optional: Docker 20.10+, Docker Compose

## Project Architecture (Quick Recap)
- DDD + layered architecture (domain, application, infrastructure, interfaces)
- Services: auth-service (HTTP:8080), gateway-service (TCP:9090), game-service (Go RPC:8081)
- Storage: MongoDB (primary), Redis (cache)
- Protocols: HTTP (auth), TCP (gateway), Go RPC (gateway↔game)
- Protobuf definitions in /proto; prefer using scripts/generate_proto.(bat|sh)

See more details in .github/instructions/gw.instructions.md.

## Getting Started
- Fork and clone the repo
- Create a branch from main: feature/<short-name> or fix/<short-name>
- Run:
  - go fmt ./...
  - go mod tidy (after dependency changes)
  - go test ./...

## Development Standards
- Follow Go naming conventions; keep package names short and lower-case.
- Use internal/infrastructure/logging for all logs (no fmt.Println).
- Always pass context.Context across boundaries; honor cancellation/timeouts.
- Keep domain entities’ fields private; expose behavior via methods.
- Application layer orchestrates use cases; do not access DB/Redis directly here.
- Infrastructure layer implements persistence/cache/messaging/network.
- Interface adapters (HTTP/TCP/RPC) map DTO/proto <-> domain models.

## Errors and Observability
- Use internal/errors for domain errors and map appropriately in interfaces.
- Wrap errors with context: fmt.Errorf("...: %w", err)
- Structured logging (json) with key fields: service, module, player_id, trace_id.
- Optional pprof per service via configs/*.yaml; expose only in trusted networks.

## Protocols & Compatibility
- When changing .proto files, preserve backward compatibility:
  - Only add fields; do not change or reuse existing tags
  - Deletions should be reserved
- Update gateway and game handling logic consistently when protocol changes
- Generate code with scripts/generate_proto.bat (Windows) or .sh (Unix)

## Configuration
- All services load from configs/*.yaml via internal/config; never hardcode ports/secrets/URIs.
- Adding new config fields requires:
  - Updating the example files under configs/*
  - Updating loading/validation and defaults
  - Documenting changes in README or service docs

## Dependencies
- Prefer stdlib and existing deps; avoid introducing heavy/CGO deps.
- If adding a new dependency:
  1) Justify the choice and alternatives; 2) Pin versions; 3) Ensure Go 1.24 compatibility.
- After changes: go mod tidy

## Testing
- Unit tests co-located with code (_test.go), table-driven where possible.
- Cover core branches; ensure `go test ./...` passes before submitting PRs.
- Integration/E2E via tools/simclient (smoke/feature/load modes). Provide minimal runnable scripts/configs for new protocols or interfaces.

## Commit Messages and PRs
- Prefer Conventional Commits:
  - feat: new feature
  - fix: bug fix
  - refactor: refactoring without behavior change
  - docs: documentation updates
  - test: testing-related changes
  - chore/build/ci: build or pipeline updates
- Keep changes small and focused; avoid unrelated refactors.
- In PR description include:
  - Motivation and design summary
  - Impact surface (protocol/config/data/perf)
  - Verification steps (unit tests and/or simclient scenarios)
  - Any config/script/docs updates included

## Running Locally
- Windows:
  - scripts/start-services.bat
- Linux/Mac:
  - ./scripts/start-services.sh
- Or run each service via `go run cmd/<service>/main.go`

## Generating Protobuf Code
- Windows: `scripts/generate_proto.bat`
- Unix: `scripts/generate_proto.sh`
- Do not manually scatter generated files; use the scripts and keep backward compatibility.

## Do / Don't Quick List
Do:
- Keep changes minimal and reversible
- Respect layering and DDD boundaries
- Log via infrastructure/logging
- Pass context and close resources via defer

Don’t:
- Break proto compatibility (tag changes or reuse)
- Hardcode ports/secrets/URIs
- Introduce heavy/CGO deps without strong justification
- Access DB/Redis from interface or application layers

## Questions
Open an issue or start a discussion if anything is unclear. Thanks for contributing!
