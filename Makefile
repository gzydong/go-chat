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
.PHONY: build-windows
build-windows:generate lint
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./windows/bin/http-server.exe ./internal/http
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./windows/bin/ws-server.exe ./internal/websocket
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./windows/bin/job-cli.exe ./internal/job
	cp ./config.example.yaml ./windows/config.yaml

.PHONY: generate
generate:
	wire ./...

lint:
	golangci-lint run --timeout=5m --config ./.golangci.yml

test:
	go test -v ./...


tool:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./script/linux-ws-tool ./script
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./script/windows-ws-tool.exe ./script