// internal/trading/trading_strategy.go

package trading

import (
	"errors"
	"fmt"
	"time"

	"go-optimal-stop/internal/data"
)

// TradingStrategy 関数の定義
func TradingStrategy(data *[]data.Data, startDate string, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (time.Time, time.Time, float64, error) {
	d := *data
	if len(d) == 0 {
		return time.Time{}, time.Time{}, 0, errors.New("データが空です")
	}

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return time.Time{}, time.Time{}, 0, fmt.Errorf("日付の解析エラー: %v", err)
	}

	maxDate := d[len(d)-1].Date
	for !dateInData(d, startDateParsed) {
		if startDateParsed.After(maxDate) {
			return time.Time{}, time.Time{}, 0, errors.New("開始日がデータの範囲外です。無限ループを防ぐため、処理を中断")
		}
		startDateParsed = startDateParsed.AddDate(0, 0, 1)
	}

	var purchaseDate time.Time
	var purchasePrice float64
	for _, day := range d {
		if day.Date.Equal(startDateParsed) {
			purchaseDate = day.Date
			purchasePrice = day.Open
			break
		}
	}

	stopLossThreshold := round(purchasePrice * (1 - stopLossPercentage/100))
	trailingStopTriggerPrice := round(purchasePrice * (1 + trailingStopTrigger/100))

	var endDate time.Time
	var endPrice float64
	for _, day := range d {
		if day.Date.Before(startDateParsed) {
			continue
		}
		openPrice := day.Open
		lowPrice := day.Low
		closePrice := day.Close

		if openPrice <= stopLossThreshold {
			endPrice = openPrice
			endDate = day.Date
			break
		}
		if lowPrice <= stopLossThreshold {
			endPrice = lowPrice
			endDate = day.Date
			break
		}
		if closePrice >= trailingStopTriggerPrice {
			stopLossThreshold = round(closePrice * (1 - trailingStopUpdate/100))
			trailingStopTriggerPrice = round(closePrice * (1 + trailingStopTrigger/100))
		}
	}
	if endDate.IsZero() {
		endPrice = d[len(d)-1].Close
		endDate = d[len(d)-1].Date
	}
	profitLoss := round((endPrice - purchasePrice) / purchasePrice * 100)
	return purchaseDate, endDate, profitLoss, nil
}

// ユーティリティ関数: 特定の日付がデータに存在するか確認
func dateInData(data []data.Data, date time.Time) bool {
	for _, day := range data {
		if day.Date.Equal(date) {
			return true
		}
	}
	return false
}

// ユーティリティ関数: 四捨五入
func round(value float64) float64 {
	return float64(int(value*10+0.5)) / 10
}
