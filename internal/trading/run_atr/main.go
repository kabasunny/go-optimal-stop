package main

import (
	"fmt"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/internal/trading" // ここに trading パッケージのインポートを追加
)

func main() {
	// テスト用のデータを作成
	dailyData := []ml_stockdata.InMLDailyData{
		{Date: "2025-01-01", High: 100.0, Low: 90.0, Close: 95.0},
		{Date: "2025-01-02", High: 102.0, Low: 91.0, Close: 96.0},
		{Date: "2025-01-03", High: 104.0, Low: 92.0, Close: 98.0},
		{Date: "2025-01-04", High: 105.0, Low: 93.0, Close: 99.0},
		{Date: "2025-01-05", High: 106.0, Low: 94.0, Close: 100.0},
		{Date: "2025-01-06", High: 107.0, Low: 95.0, Close: 101.0},
		{Date: "2025-01-07", High: 108.0, Low: 96.0, Close: 102.0},
		{Date: "2025-01-08", High: 109.0, Low: 97.0, Close: 103.0},
		{Date: "2025-01-09", High: 110.0, Low: 98.0, Close: 104.0},
		{Date: "2025-01-10", High: 111.0, Low: 99.0, Close: 105.0},
		{Date: "2025-01-11", High: 112.0, Low: 100.0, Close: 106.0},
		{Date: "2025-01-12", High: 113.0, Low: 101.0, Close: 107.0},
		{Date: "2025-01-13", High: 114.0, Low: 102.0, Close: 108.0},
		{Date: "2025-01-14", High: 120.0, Low: 110.0, Close: 115.0},
	}

	signalDate, _ := time.Parse("2006-01-02", "2025-01-14")
	StopLossPercentage := 2.0
	portfolioValue := 1000000
	availableFundsInt := 500000
	entryPrice := 115.0
	commissionRate := 0.1

	positionSize, entryCost, err := trading.DeterminePositionSize(StopLossPercentage, portfolioValue, availableFundsInt, entryPrice, &commissionRate, &dailyData, signalDate)

	if err != nil {
		fmt.Printf("エラーが発生しました: %v\n", err)
	} else {
		fmt.Printf("ポジションサイズ: %.2f株\n", positionSize)
		fmt.Printf("エントリーコスト: %.2f円\n", entryCost)
	}
}
