LDFLAGS="-X main.commit=$(shell git rev-parse --short HEAD) -X main.buildDate=$(shell date +"%Y-%m-%dT%H:%M:%S%z")"

.PHONY: build
build:
	CGO_ENABLED=1 go build -ldflags=$(LDFLAGS) -o ./bin/relay ./cmd/relay

.PHONY: run
run: build
	./bin/relay -config local/config.yml

.PHONY: test
test:
	go test -v ./...

.PHONY: run-docker
run-docker:
	@docker compose up

.PHONY: clean-docker
clean-docker:
	docker compose kill || true
	docker compose rm -f || true

.PHONY: build-docker
build-docker:
	docker compose build

