# Go 1.23.2 버전을 사용하여 빌드
FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY . .

# go.mod 파일이 없으면 생성하고 의존성 정리
RUN go mod init rsi-bot || true
RUN go mod tidy
RUN go build -o rsi-bot main.go

# 실행을 위한 가벼운 이미지
FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache tzdata
ENV TZ=Asia/Seoul
COPY --from=builder /app/rsi-bot .

# 실행 권한 부여 및 실행
RUN chmod +x rsi-bot
CMD ["./rsi-bot"]
