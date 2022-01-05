conf:
	cp config.example.yaml config.yaml

.PHONY: run
run: generate
	go run .

.PHONY: build
build:generate lint
	go build -o ./bin/app

.PHONY: generate
generate:
	wire ./...

websocket:generate
	go run ./app/websocket/

lint:
	golangci-lint run --timeout=5m --config ./.golangci.yml

test:
	go test -v ./...