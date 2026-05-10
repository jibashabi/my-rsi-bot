# Go 1.23.2 버전을 사용하여 빌드
FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY . .
# 의존성 초기화 및 정리
RUN go mod init rsi-bot || true
RUN go mod tidy
RUN go build -o rsi-bot main.go

FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache tzdata
ENV TZ=Asia/Seoul
COPY --from=builder /app/rsi-bot .
CMD ["./rsi-bot"]
