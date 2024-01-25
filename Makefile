PROTO_FILES := $(shell find api -name *.proto)
OS:=$(shell go env GOHOSTOS)

ifeq ($(OS), windows)
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	PROTO_FILES=$(shell $(Git_Bash) -c "find api -iname *.proto")
else
	PROTO_FILES=$(shell find api -iname *.proto)
endif

.PHONY: install
install:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/srikrsna/protoc-gen-gotag@latest

.PHONY: conf
conf:
	cp config.example.yaml config.yaml

.PHONY: generate
generate:
	go generate ./...

lint:
	golangci-lint run --timeout=5m --config ./.golangci.yml

test:
	go test -v ./...

http:
	go run ./cmd/lumenim http

comet:
	go run ./cmd/lumenim comet

migrate:
	go run ./cmd/lumenim migrate

queue:
	go run ./cmd/lumenim queue

crontab:
	go run ./cmd/lumenim crontab

.PHONY: build
build:
	go build -o ./bin/lumenim ./cmd/lumenim

.PHONY: build-all
build-all:
	@mkdir -p ./build/linux/ ./build/windows/ ./build/mac/ ./build/macm1/

	# 构建 windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/windows/lumenim ./cmd/lumenim
	cp ./config.example.yaml ./build/windows/config.yaml

	# 构建 linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/linux/lumenim ./cmd/lumenim
	cp ./config.example.yaml ./build/linux/config.yaml

	# 构建 mac amd
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/mac/lumenim ./cmd/lumenim
	cp ./config.example.yaml ./build/mac/config.yaml

	# 构建 mac m1
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./build/macm1/lumenim ./cmd/lumenim
	cp ./config.example.yaml ./build/macm1/config.yaml

.PHONY: proto
proto:
	@if [ -n "$(PROTO_FILES)" ]; then \
		protoc \
		--proto_path=./api/proto \
		--proto_path=./third_party \
		--go_out=paths=source_relative:./api/pb/ \
		--validate_out=paths=source_relative,lang=go:./api/pb/ $(PROTO_FILES) \
	 && protoc --proto_path=./third_party --proto_path=./api/proto --gotag_out=outdir="./api/pb/":./ $(PROTO_FILES) \
	 && echo "protoc generate success"; \
	fi


.PHONY: deploy
deploy:
	git reset --hard origin/develop && git pull && make build && supervisorctl reload

#--go-grpc_out=paths=source_relative:./api/pb/ \
