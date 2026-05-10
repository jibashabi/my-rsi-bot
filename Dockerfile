# 1. 넉넉하게 최신 버전인 1.23.2을 사용합니다.
FROM golang:1.23.2-alpine AS builder

# 2. 필요한 빌드 도구들을 설치합니다 (alpine은 가벼워서 직접 설치가 필요할 때가 있습니다).
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app
COPY . .

# 3. 의존성 파일을 강제로 초기화하고 필요한 라이브러리를 가져옵니다.
RUN go mod init rsi-bot || true
RUN go mod tidy
RUN go build -o rsi-bot main.go

# 4. 실제 실행할 가벼운 환경으로 복사합니다.
FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache tzdata
ENV TZ=Asia/Seoul

COPY --from=builder /app/rsi-bot .

# 5. 실행 권한을 주고 봇을 가동합니다.
RUN chmod +x rsi-bot
CMD ["./rsi-bot"]
