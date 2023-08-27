FROM golang:1.19-alpine as builder

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /home/admin
COPY . .

RUN go mod download

# CGO_ENABLED禁用cgo 然后指定OS等，并go build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chaosmeta ./cmd/server

# 运行阶段指定alpine作为基础镜像
FROM alpine:3.12

WORKDIR /home/admin

COPY --from=builder /home/admin/chaosmeta .

EXPOSE 8080

#ADD /conf/app.yaml /home/admin
#ADD /conf/app.conf /home/admin

ENTRYPOINT  ["./chaosmeta", "server", "start", "--config=conf/app.yaml"]
