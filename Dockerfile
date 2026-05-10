# 1. 빌드 도구 자체를 1.23.2 버전으로 강제 지정
FROM golang:1.23.2-bookworm AS builder

# 2. 환경 변수로 Go 버전 체크를 무시하도록 설정
ENV GO111MODULE=on
ENV GOTOOLCHAIN=go1.23.2

WORKDIR /app
COPY . .

# 3. 기존 설정 무시하고 새로 구성
RUN rm -f go.mod go.sum || true
RUN go mod init rsi-bot
RUN go mod tidy

RUN CGO_ENABLED=0 go build -o rsi-bot main.go

# 4. 실행 환경 (더 안정적인 환경 사용)
FROM debian:bookworm-slim
WORKDIR /root/
RUN apt-get update && apt-get install -y ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
ENV TZ=Asia/Seoul
COPY --from=builder /app/rsi-bot .

RUN chmod +x rsi-bot

# 아래 한 줄을 추가해 주세요!
EXPOSE 8080

CMD ["./rsi-bot"]
