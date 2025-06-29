# 多阶段构建 Dockerfile
# 第一阶段：构建前端
FROM node:20-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app/frontend

# 复制前端依赖文件
COPY WEBSITE/package*.json ./

# 安装依赖
RUN npm install

# 复制前端源代码
COPY WEBSITE/ ./

# 构建前端项目
RUN npm run build

# 第二阶段：构建后端
FROM golang:1.23-alpine AS backend-builder

# 安装必要的工具和sqlite开发包
# RUN apk add --no-cache go

# 设置工作目录
WORKDIR /app/backend

# 复制Go模块文件
COPY GO/go.mod GO/go.sum ./

# 下载依赖
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
RUN go mod download

# 复制后端源代码
COPY GO/ ./

# 构建后端可执行文件
ENV CGO_ENABLED=0
RUN go build -a -installsuffix cgo -o main main.go

# 第三阶段：最终运行镜像
FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# 安装运行时依赖，包括sqlite
RUN apk --no-cache add ca-certificates tzdata sqlite

# 设置时区
ENV TZ=Asia/Shanghai

# 创建工作目录
WORKDIR /app

# 从构建阶段复制前端构建产物
COPY --from=frontend-builder /app/frontend/dist ./dist

# 从构建阶段复制后端可执行文件
COPY --from=backend-builder /app/backend/main ./

# 复制配置文件
COPY GO/.env ./

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 创建数据目录并设置权限
RUN mkdir -p ./data && chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 设置默认端口环境变量
ENV PORT=8081

# 暴露端口
EXPOSE $PORT

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:$PORT/health || exit 1

# 启动应用
CMD ["./main"]
