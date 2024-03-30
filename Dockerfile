## 容器交叉编译, 并上传
#FROM golang as builder
#
#ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
#    GO111MODULE=on
#
##设置时区参数
#ENV TZ=Asia/Shanghai
#
#WORKDIR /app
#COPY . /app
#
## 编译应用程序为x86架构
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" .
#
#
#
## 第二阶段：使用scratch作为基础镜像创建最终镜像
#FROM --platform=linux/amd64 Alpine
#
## 从构建器阶段复制编译好的应用程序
#COPY --from=builder /app/eduData /eduData
#
## 运行应用程序
#ENTRYPOINT ["./eduData"]

#FROM alpine

FROM alpine:3

# 待解决时区问题
ENV TZ=Asia/Shanghai
#RUN apk add tzdata

COPY ./eduData /
COPY config/config.ini /config/config.ini

EXPOSE 8080

ENTRYPOINT ["/eduData"]