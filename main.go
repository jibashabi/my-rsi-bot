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

	fmt.Println("🚀 RSI 실시간 감시 시작 (OKX)")

	for {
		resp, err := client.R().Get("https://www.okx.com/api/v5/market/ticker?instId=BTC-USDT")
		if err != nil {
			log.Println("API 요청 실패:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		var result OKXResponse
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			log.Println("데이터 변환 실패:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if len(result.Data) > 0 {
			price, _ := strconv.ParseFloat(result.Data[0].Last, 64)
			prices = append(prices, price)

			// 데이터가 14개 이상 쌓이면 RSI 계산
			if len(prices) > 14 {
				rsi := talib.Rsi(prices, 14)
				currentRSI := rsi[len(rsi)-1]
				fmt.Printf("[%s] 현재가: %.2f | RSI: %.2f\n", time.Now().Format("15:04:05"), price, currentRSI)

				// 오래된 데이터 삭제 (메모리 관리)
				prices = prices[1:]
			} else {
				fmt.Printf("[%s] 데이터 수집 중... (%d/15)\n", time.Now().Format("15:04:05"), len(prices))
			}
		}

		time.Sleep(10 * time.Second) // 10초마다 확인
	}
}
