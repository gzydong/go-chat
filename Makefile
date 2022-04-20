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
	go build -o ./bin/http ./internal/http
	go build -o ./bin/websocket ./internal/websocket
	go build -o ./bin/job-cli ./internal/job

## mac 下打包 windows 执行文件
.PHONY: build-all
build-all:generate lint build-windows build-linux build-mac

.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/windows/bin/http-server.exe ./internal/http
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/windows/bin/ws-server.exe ./internal/websocket
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/windows/bin/job-cli.exe ./internal/job
	cp ./config.example.yaml ./build/windows/config.yaml

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/linux/bin/http-server ./internal/http
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/linux/bin/ws-server ./internal/websocket
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/linux/bin/job-cli ./internal/job
	cp ./config.example.yaml ./build/linux/config.yaml

.PHONY: build-mac
build-mac:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/mac/bin/http-server ./internal/http
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/mac/bin/ws-server ./internal/websocket
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./build/mac/bin/job-cli ./internal/job
	cp ./config.example.yaml ./build/mac/config.yaml

.PHONY: generate
generate:
	wire ./...

lint:
	golangci-lint run --timeout=5m --config ./.golangci.yml

test:
	go test -v ./...


tool:
	go build -o ./script/mac-ws-tool ./script
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./script/linux-ws-tool ./script
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./script/windows-ws-tool.exe ./script