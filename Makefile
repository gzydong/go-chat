PROTO_FILES := $(shell find api -iname *.proto)

.PHONY: install
install:
	go install github.com/google/wire/cmd/wire@latest \
	&& go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
	&& go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
	&& go install github.com/envoyproxy/protoc-gen-validate@latest \
	&& go install github.com/srikrsna/protoc-gen-gotag@latest \

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

http: generate
	go run ./internal/http

im_server: generate
	go run ./internal/gateway

cmd: generate
	go run ./internal/cmd

migrate:
	go run ./internal/cmd other migrate

.PHONY: build
build:generate
	go build -o ./bin/http-server ./internal/http
	go build -o ./bin/im-server ./internal/gateway
	go build -o ./bin/cmd ./internal/cmd

.PHONY: build-all
build-all:generate lint
	# 构建 windows
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/windows/bin/http-server.exe ./internal/http
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/windows/bin/im-server.exe ./internal/gateway
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/windows/bin/cmd.exe ./internal/cmd
	cp ./config.example.yaml ./build/windows/config.yaml

	# 构建 linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/linux/bin/http-server ./internal/http
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/linux/bin/im-server ./internal/gateway
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/linux/bin/cmd ./internal/cmd
	cp ./config.example.yaml ./build/linux/config.yaml

	# 构建 mac amd
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/mac/bin/http-server ./internal/http
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/mac/bin/im-server ./internal/gateway
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./build/mac/bin/cmd ./internal/cmd
	cp ./config.example.yaml ./build/mac/config.yaml

	# 构建 mac m1
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./build/macm1/bin/http-server ./internal/http
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./build/macm1/bin/im-server ./internal/gateway
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./build/macm1/bin/cmd ./internal/cmd
	cp ./config.example.yaml ./build/macm1/config.yaml

.PHONY: proto
proto:
	@if [ -n "$(PROTO_FILES)" ]; then \
		protoc --proto_path=./api/proto \
		--proto_path=./third_party \
		--go_out=paths=source_relative:./api/pb/ \
		--validate_out=paths=source_relative,lang=go:./api/pb/ $(PROTO_FILES) \
	 && protoc --proto_path=./third_party --proto_path=./api/proto --gotag_out=outdir="./api/pb/":./ $(PROTO_FILES) \
	 && echo "protoc generate success"; \
	fi


.PHONY: deploy
deploy:
	git pull && make build && supervisorctl reload

#--go-grpc_out=paths=source_relative:./api/pb/ \
