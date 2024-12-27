// internal/trading/stop_loss_utils.go

package trading

import (
	"go-optimal-stop/internal/ml_stockdata"
	"time"
)

// calculateStopLoss 関数: 損切りしきい値とトリガー価格を計算
func calculateStopLoss(purchasePrice, stopLossPercentage, trailingStopTrigger float64) (float64, float64) {
	stopLossThreshold := round(purchasePrice * (1 - stopLossPercentage/100))
	trailingStopTriggerPrice := round(purchasePrice * (1 + trailingStopTrigger/100))
	return stopLossThreshold, trailingStopTriggerPrice
}

// findExitDate 関数: 退出日を見つける
func findExitDate(data []ml_stockdata.MLDailyData, startDate time.Time, stopLossThreshold, trailingStopTriggerPrice, trailingStopTrigger, trailingStopUpdate float64) (time.Time, float64, error) {
	var endDate time.Time
	var endPrice float64
	for _, day := range data {
		parsedDate, err := parseDate(day.Date)
		if err != nil {
			return time.Time{}, 0, err
		}
		if parsedDate.Before(startDate) {
			continue
		}
		openPrice := day.Open
		lowPrice := day.Low
		closePrice := day.Close

		if openPrice <= stopLossThreshold {
			endPrice = openPrice
			endDate = parsedDate
			break
		}
		if lowPrice <= stopLossThreshold {
			endPrice = lowPrice
			endDate = parsedDate
			break
		}
		if closePrice >= trailingStopTriggerPrice {
			stopLossThreshold = round(closePrice * (1 - trailingStopUpdate/100))
			trailingStopTriggerPrice = round(closePrice * (1 + trailingStopTrigger/100))
		}
	}
	if endDate.IsZero() {
		endPrice = data[len(data)-1].Close
		endDate, _ = parseDate(data[len(data)-1].Date)
	}
	return endDate, endPrice, nil
}

// round 関数: 四捨五入
func round(value float64) float64 {
	return float64(int(value*10+0.5)) / 10
}
