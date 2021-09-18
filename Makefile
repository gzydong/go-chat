.PHONY: conf
conf:
	cp config.yaml.example config.yaml

run:
	go run .

.PHONY: build
build:
	go build -o go-chat