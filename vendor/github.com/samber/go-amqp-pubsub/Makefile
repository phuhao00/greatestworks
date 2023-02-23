
all: build
re: clean all

#
# Build
#
build:
	go build -v ./example/producer/*.go
watch-build: deps
	reflex -t 50ms -s -- sh -c 'echo \\nBUILDING && CGO_ENABLED=0 dlv --listen=:1234 --headless=true --accept-multiclient --api-version=2 debug ./example/producer/*.go --continue && echo Exited'

#
# Deps
#
deps-tools:
	go install github.com/cespare/reflex@latest
	go install github.com/rakyll/gotest@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install github.com/psampaz/go-mod-outdated@latest
	go install github.com/jondot/goweight@latest
	go install golang.org/x/tools/cmd/cover@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/sonatype-nexus-community/nancy@latest
	go mod tidy

deps: deps-tools
	go mod download -x

cleanup-deps:
	go mod tidy

audit: cleanup-deps
	go list -json -m all | nancy sleuth

outdated: cleanup-deps
	go list -u -m -json all | go-mod-outdated -update -direct

vulncheck:
	govulncheck ./...

#
# Quality
#
lint:
	golangci-lint run --timeout 600s --max-same-issues 50 --path-prefix=./ ./...
lint-fix:
	golangci-lint run --fix --timeout 600s --max-same-issues 50 --path-prefix=./ ./...

test:
	gotest -v ./...
watch-test: deps
	reflex -t 50ms -s -- sh -c 'make test'

weight:
	goweight

coverage:
	go test -v -coverprofile=cover.out -covermode=atomic ./...
	go tool cover -html=cover.out -o cover.html
