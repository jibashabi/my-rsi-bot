package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/markcheno/go-talib"
)

const (
	NTFY_URL       = "https://ntfy.sh/rsi111"
	CHECK_INTERVAL = 1 * time.Minute // 1분마다 실시간 체크
	RSI_PERIOD     = 14
	RSI_OVERBOUGHT = 70.0
	RSI_OVERSOLD   = 30.0
)

// 중복 알림 방지를 위한 상태 저장 (메모리상)
var lastStatus = make(map[string]string) 

func main() {
	symbols := []string{"BTC-USDT", "ETH-USDT"}
	client := resty.New()

	fmt.Println("⚡ RSI 실시간 감시 시작 (1시간 차트 기준)")

	for {
		for _, symbol := range symbols {
			checkRealtimeRSI(client, symbol)
		}
		time.Sleep(CHECK_INTERVAL)
	}
}

func checkRealtimeRSI(client *resty.Client, symbol string) {
	// OKX에서 1시간 봉 데이터 가져오기 (가장 최근 100개)
	resp, err := client.R().
		SetQueryParam("instId", symbol).
		SetQueryParam("bar", "1H"). 
		SetQueryParam("limit", "100").
		Get("https://www.okx.com/api/v5/market/candles")

	if err != nil {
		return
	}

	// ... (데이터 파싱 로직은 이전과 동일, 생략) ...
    // 여기서 Data[0]은 현재 실시간으로 움직이는 '1시간 봉'의 데이터입니다.

	var closes []float64
	// OKX 데이터를 과거 순으로 정렬 (최신 봉이 마지막에 오도록)
	// okxResp.Data[0]이 현재가(실시간)임
	for i := 99; i >= 0; i-- {
		val, _ := strconv.ParseFloat(okxResp.Data[i][4], 64)
		closes = append(closes, val)
	}

	// 현재 실시간 RSI 계산
	rsi := talib.Rsi(closes, RSI_PERIOD)
	currentRSI := rsi[len(rsi)-1]

	fmt.Printf("[%s] 실시간 RSI: %.2f\n", symbol, currentRSI)

	// 알림 조건 판단 및 중복 방지
	if currentRSI <= RSI_OVERSOLD && lastStatus[symbol] != "OVERSOLD" {
		sendNotification(client, symbol, "🔴 과매도 진입!", currentRSI)
		lastStatus[symbol] = "OVERSOLD"
	} else if currentRSI >= RSI_OVERBOUGHT && lastStatus[symbol] != "OVERBOUGHT" {
		sendNotification(client, symbol, "🟢 과매수 진입!", currentRSI)
		lastStatus[symbol] = "OVERBOUGHT"
	} else if currentRSI > RSI_OVERSOLD && currentRSI < RSI_OVERBOUGHT {
		// 정상 범위로 돌아오면 상태 초기화
		lastStatus[symbol] = "NORMAL"
	}
}

func sendNotification(client *resty.Client, symbol, status string, rsi float64) {
	msg := fmt.Sprintf("[%s] %s (RSI: %.2f)", symbol, status, rsi)
	client.R().SetBody(msg).Post(NTFY_URL)
	fmt.Println("📢 알림 발송:", msg)
}
