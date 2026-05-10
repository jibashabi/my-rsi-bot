FROM golang:1.22.1-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod init rsi-bot && go mod tidy
RUN go build -o rsi-bot main.go

FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache tzdata
ENV TZ=Asia/Seoul
COPY --from=builder /app/rsi-bot .
CMD ["./rsi-bot"]
