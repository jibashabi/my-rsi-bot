package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/markcheno/go-talib"
)

type OKXResponse struct {
	Code string `json:"code"`
	Data []struct {
		Last string `json:"last"`
	} `json:"data"`
}

func main() {
	client := resty.New()
	var prices []float64
	// 사용자님의 ntfy 주소
	ntfyURL := "https://ntfy.sh/rsi111" 

	fmt.Println("🚀 RSI 알림 봇 가동 시작 (채널: rsi111)")

	// [테스트용] 봇이 실행되자마자 알림을 한 번 보냅니다.
	client.R().SetBody("✅ RSI 알림 봇이 정상적으로 연결되었습니다! (테스트)").Post(ntfyURL)

	for {
		resp, err := client.R().Get("https://www.okx.com/api/v5/market/ticker?instId=BTC-USDT")
		if err != nil {
			log.Println("API 요청 실패:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		var result OKXResponse
		json.Unmarshal(resp.Body(), &result)

		if len(result.Data) > 0 {
			price, _ := strconv.ParseFloat(result.Data[0].Last, 64)
			prices = append(prices, price)

			if len(prices) > 14 {
				rsi := talib.Rsi(prices, 14)
				currentRSI := rsi[len(rsi)-1]
				fmt.Printf("[%s] 가격: %.2f | RSI: %.2f\n", time.Now().Format("15:04:05"), price, currentRSI)

				// 실제 알림 조건
				if currentRSI >= 70 {
					client.R().SetBody(fmt.Sprintf("🔥 과매수! BTC RSI: %.2f", currentRSI)).Post(ntfyURL)
				} else if currentRSI <= 30 {
					client.R().SetBody(fmt.Sprintf("🧊 과매도! BTC RSI: %.2f", currentRSI)).Post(ntfyURL)
				}
				
				prices = prices[1:]
			} else {
				fmt.Printf("데이터 수집 중... (%d/15)\n", len(prices))
			}
		}
		time.Sleep(10 * time.Second)
	}
}
