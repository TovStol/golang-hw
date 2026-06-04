BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_STORER := "./bin/calendar_storer"
DOCKER_IMG="calendar:develop"
GO ?= go

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	$(GO) build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	$(GO) build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	$(GO) build -v -o $(BIN_STORER) -ldflags "$(LDFLAGS)" ./cmd/calendar_storer

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	$(GO) test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.2

lint: install-lint-deps
	golangci-lint run ./...

up:
	cd deployments && docker compose up -d

down:
	cd deployments && docker compose down

integration-tests:
	cd deployments && docker compose -f docker-compose.test.yaml down -v || true
	cd deployments && docker compose -f docker-compose.test.yaml build
	cd deployments && docker compose -f docker-compose.test.yaml up -d
	sleep 45
	cd deployments && docker compose -f docker-compose.test.yaml exec integration-tests go test -tags integration -v -timeout 120s ./integration/...
	cd deployments && docker compose -f docker-compose.test.yaml down -v

.PHONY: build run build-img run-img version test lint run-scheduler run-storer up down integration-tests
