.PHONY: conf
conf:
	cp config.example.yaml config.yaml

run:
	go run .

.PHONY: build
build:
	go build -o go-chat