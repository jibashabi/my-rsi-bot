# 1. 넉넉하게 최신 버전인 1.23.2을 사용합니다.
FROM golang:1.23.2-alpine AS builder

# 2. 빌드에 필요한 도구들을 미리 설치합니다.
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app
COPY . .

# 3. 모듈 설정 파일(go.mod)을 강제로 초기화하고 라이브러리를 다시 받습니다.
RUN rm -f go.mod go.sum || true
RUN go mod init rsi-bot
RUN go mod tidy

# 4. 빌드 진행 (CGO를 꺼서 호환성을 높입니다)
RUN CGO_ENABLED=0 go build -o rsi-bot main.go

# 5. 실행 환경 (실제 봇이 돌아가는 곳)
FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache tzdata
ENV TZ=Asia/Seoul

COPY --from=builder /app/rsi-bot .

# 6. 실행 권한 부여 및 시작
RUN chmod +x rsi-bot
CMD ["./rsi-bot"]
