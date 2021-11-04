.PHONY: conf
conf:
	cp config.example.yaml config.yaml

.PHONY: run
run: generate
	go run .

.PHONY: build
build:generate
	go build -o ./bin/app

.PHONY: generate
generate:
	wire

.PHONY: lint
lint:
	golangci-lint run ./...