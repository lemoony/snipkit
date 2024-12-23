.DEFAULT_GOAL := dev

.PHONY: dev
dev: ## dev build
dev: clean install generate vet fmt lint test mod-tidy

.PHONY: ci
ci: ## CI build
ci: dev diff

.PHONY: clean
clean: ## remove files created during build pipeline
	$(call print-target)
	rm -rf dist
	rm -f coverage.*

.PHONY: install
install: ## go install tools
	$(call print-target)
	cd tools && go install $(shell cd tools && go list -e -f '{{ join .Imports " " }}' -tags=tools)

.PHONY: generate
generate: mocks ## go generate
	$(call print-target)
	go generate ./...

.PHONY: mocks
mocks: ## go generate
	$(call print-target)
	mockery --name='App' --output="./mocks/app" --dir="./internal/app"
	mockery --name='Assistant' --output="./mocks/assistant" --dir="./internal/assistant"
	mockery --name='Client' --output="./mocks/assistant/client" --dir="./internal/assistant"
	mockery --name='Service' --output="./mocks/config" --structname="ConfigService"  --dir="./internal/config"
	mockery --name='Provider' --output="./mocks/managers"  --dir="./internal/managers"
	mockery --name='Manager' --output="./mocks/managers" --dir="./internal/managers"
	mockery --name='Cache' --output="./mocks/cache" --dir="./internal/cache"
	mockery --name='^TUI$$' --output="./mocks/ui" --dir="./internal/ui"
	mockery --name='MessagePrinter' --output="./mocks/ui" --dir="./internal/ui"
	mockery --name='Screen' --output="./mocks/ui/sync" --structname="SyncScreen" --dir="./internal/ui/sync"

.PHONY: vet
vet: ## go vet
	$(call print-target)
	go vet ./...

.PHONY: fmt
fmt: ## go fmt
	$(call print-target)
	go fmt ./...
	gci write -s standard -s default -s "prefix(github.com/lemoony/snipkit)" main.go ./cmd ./internal ./themes
	gofumpt -l -w .

.PHONY: lint
lint: ## golangci-lint
	$(call print-target)
	golangci-lint run

.PHONY: test
test: ## go test with race detector and code covarage
	$(call print-target)
	go test -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: mod-tidy
mod-tidy: ## go mod tidy
	$(call print-target)
	go mod tidy
	cd tools && go mod tidy

.PHONY: diff
diff: ## git diff
	$(call print-target)
	git diff --exit-code
	RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi

.PHONY: build
build: ## goreleaser --snapshot --skip=publish --clean
build: install
	$(call print-target)
	goreleaser --snapshot --skip=publish --clean

.PHONY: build-demo
build-demo: ## BUILD_TAGS="demo" goreleaser --snapshot --skip=publish --clean
build-demo: install
	$(call print-target)
	BUILD_TAGS="demo" goreleaser --snapshot --skip=publish --clean

.PHONY: release
release: ## goreleaser --rm-dist
release: install
	$(call print-target)
	goreleaser --rm-dist

.PHONY: run
run: ## go run
	@go run -race .

.PHONY: go-clean
go-clean: ## go clean build, test and modules caches
	$(call print-target)
	go clean -r -i -cache -testcache -modcache

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
