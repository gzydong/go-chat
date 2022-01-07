conf:
	cp config.example.yaml config.yaml

http: generate
	go run ./internal/http

websocket: generate
	go run ./internal/websocket
job: generate
	go run ./internal/job

.PHONY: build
build:generate lint
	go build -o ./bin/http ./internal/http/
	go build -o ./bin/websocket ./internal/websocket/
	go build -o ./bin/job ./internal/job/

.PHONY: generate
generate:
	wire ./...

lint:
	golangci-lint run --timeout=5m --config ./.golangci.yml

test:
	go test -v ./...