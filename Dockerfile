FROM golang as builder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=1

#设置时区参数
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY . /app

RUN go build .

EXPOSE 8080

ENTRYPOINT ["./eduData"]