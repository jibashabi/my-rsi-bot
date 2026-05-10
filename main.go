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
	// 주소가 정확한지 확인해 주세요 (ntfy.sh/rsi111)
	ntfyURL := "https://ntfy.sh/rsi111" 

	fmt.Println("🚀 RSI 알림 봇 가동 시작 (채널: rsi111)")

	// [중요] 실행되자마자 무조건 알림을 쏘는 테스트 코드
	resp, err := client.R().
		SetBody("🔔 봇 가동 성공! 이제부터 RSI를 감시합니다.").
		Post(ntfyURL)

	if err != nil {
		fmt.Printf("❌ 알림 전송 에러: %v\n", err)
	} else {
		fmt.Printf("✅ 테스트 알림 전송 시도 완료 (상태코드: %d)\n", resp.StatusCode())
	}

	for {
		// OKX에서 비트코인 가격 가져오기
		apiResp, apiErr := client.R().Get("https://www.okx.com/api/v5/market/ticker?instId=BTC-USDT")
		if apiErr != nil {
			log.Println("API 요청 실패:", apiErr)
			time.Sleep(10 * time.Second)
			continue
		}

		var result OKXResponse
		if err := json.Unmarshal(apiResp.Body(), &result); err != nil {
			log.Println("데이터 변환 실패:", err)
			continue
		}

		if len(result.Data) > 0 {
			price, _ := strconv.ParseFloat(result.Data[0].Last, 64)
			prices = append(prices, price)

			// 데이터가 14개 이상 쌓였을 때만 RSI 계산
			if len(prices) > 14 {
				rsi := talib.Rsi(prices, 14)
				currentRSI := rsi[len(rsi)-1]
				fmt.Printf("[%s] 가격: %.2f | RSI: %.2f\n", time.Now().Format("15:04:05"), price, currentRSI)

				// RSI 조건 알림 (실제 상황)
				if currentRSI >= 70 {
					client.R().SetBody(fmt.Sprintf("🔥 과매수 알림! BTC RSI: %.2f", currentRSI)).Post(ntfyURL)
				} else if currentRSI <= 30 {
					client.R().SetBody(fmt.Sprintf("🧊 과매도 알림! BTC RSI: %.2f", currentRSI)).Post(ntfyURL)
				}
				
				// 메모리 관리: 가장 오래된 데이터 삭제
				prices = prices[1:]
			} else {
				fmt.Printf("데이터 수집 중... (%d/15)\n", len(prices))
			}
		}

		time.Sleep(10 * time.Second) // 10초 간격 무한 반복
	}
}
