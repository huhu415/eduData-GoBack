# 第一阶段：构建阶段
FROM registry.cn-wulanchabu.aliyuncs.com/zzyan/golang:1.23.0-alpine3.19 AS builder

# 设置时区
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata

# 设置工作目录
WORKDIR /app

# 复制源代码
COPY . .
RUN apk add --no-cache git make

# 构建应用
RUN make build

# 第二阶段：运行阶段
FROM registry.cn-wulanchabu.aliyuncs.com/zzyan/alpine:latest

# 设置时区
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/eduData /

EXPOSE 8080

HEALTHCHECK --interval=60s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/eduData"]
