# 使用官方的 Go 镜像作为构建环境
FROM golang:1.23.1-alpine AS builder

# 设置工作目录
WORKDIR /builder

# 将 go.mod 和 go.sum 文件复制到工作目录
COPY go.mod go.sum ./

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"

# 下载依赖包
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 构建可执行文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -p 4 -o lumenim ./cmd/lumenim

# 使用一个更小的基础镜像来减小最终镜像的大小
FROM alpine:latest

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 设置工作目录
WORKDIR /work/

# 从构建阶段复制可执行文件到最终镜像
COPY --from=builder /builder/lumenim .

# 运行可执行文件
ENTRYPOINT ["./lumenim"]
CMD ["--help"]