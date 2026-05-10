# 1. 최신 Go 버전 사용
FROM golang:1.23-alpine AS builder

# 2. 빌드에 필요한 도구 설치
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app
COPY . .

# 3. 모듈 파일 강제 생성 및 최신화
RUN go mod init rsi-bot || true
RUN go mod tidy

# 4. 빌드 (에러 방지를 위해 CGO_ENABLED=0 설정)
RUN CGO_ENABLED=0 go build -o rsi-bot main.go

# 5. 실행 환경 구성
FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache tzdata
ENV TZ=Asia/Seoul
COPY --from=builder /app/rsi-bot .

RUN chmod +x rsi-bot
CMD ["./rsi-bot"]
