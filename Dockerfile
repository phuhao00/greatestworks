# =============================================================================
# 多阶段构建优化版 Dockerfile
# =============================================================================

# 构建阶段 - 使用官方 Go 镜像
FROM golang:1.24-alpine AS builder

# 设置构建参数
ARG BUILD_VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT
# 选择要构建的服务包，默认构建 game-service，可传入 ./cmd/auth-service 或 ./cmd/gateway-service
ARG SERVICE_PACKAGE=./cmd/game-service

# 安装构建依赖
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx

# 设置工作目录
WORKDIR /build

# 优化 Go 模块缓存
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制源代码
COPY . .

# 构建优化的二进制文件
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -a \
    -installsuffix cgo \
    -ldflags="-s -w -X main.version=${BUILD_VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o server \
    ${SERVICE_PACKAGE}

# 使用 UPX 压缩二进制文件（可选）
RUN upx --best --lzma server

# =============================================================================
# 运行阶段 - 使用 scratch 最小化镜像
# =============================================================================
FROM scratch AS runtime

# 从构建阶段复制必要文件
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/server /server

# 复制配置文件（如果存在）
COPY --from=builder /build/configs/ /configs/

# 设置环境变量
ENV TZ=Asia/Shanghai
ENV GIN_MODE=release
ENV LOG_LEVEL=info

# 暴露端口
EXPOSE 8080 8081 9090

# 添加健康检查用户
USER 65534:65534

# 运行应用
ENTRYPOINT ["/server"]

# =============================================================================
# 开发阶段 - 包含调试工具
# =============================================================================
FROM alpine:latest AS development

# 安装开发工具
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    wget \
    netcat-openbsd \
    htop \
    strace

# 创建非root用户
RUN addgroup -g 1000 appgroup && \
    adduser -D -s /bin/sh -u 1000 -G appgroup appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制文件
COPY --from=builder /build/server .
COPY --from=builder /build/configs/ ./configs/

# 创建必要目录
RUN mkdir -p /var/log/mmo-server && \
    chown -R appuser:appgroup /app /var/log/mmo-server

# 切换到非root用户
USER appuser

# 设置环境变量
ENV TZ=Asia/Shanghai
ENV GIN_MODE=debug
ENV LOG_LEVEL=debug

# 暴露端口
EXPOSE 8080 8081 9090

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider --timeout=5 http://localhost:8080/health || exit 1

# 运行应用
CMD ["./server"]

# =============================================================================
# 默认目标为生产环境
# =============================================================================
FROM runtime AS final