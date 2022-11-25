BIN := "./bin/imgproxy"
DOCKER_IMG="imgproxy:develop"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/imgproxy

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-local: build
	$(BIN) -config ./configs/config.yaml

run: build-img
	docker-compose -f docker/docker-compose.yaml up -d

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

coverage:
	go test -race -covermode=atomic -coverprofile=coverage.out ./internal/...

intgr: build-img
	docker-compose -f ./docker/docker-compose.yaml -f ./docker/docker-compose.tests.yaml up --build --force-recreate --abort-on-container-exit --exit-code-from tests && \
	docker-compose -f ./docker/docker-compose.yaml -f ./docker/docker-compose.tests.yaml down

.PHONY: build run-local run build-img run-img version test lint coverage intgr
