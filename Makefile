.PHONY: conf
conf:
	cp config.example.yaml config.yaml

.PHONY: run
run: generate
	go run .

.PHONY: build
build:generate
	go build -o go-chat

.PHONY: generate
generate:
	wire
