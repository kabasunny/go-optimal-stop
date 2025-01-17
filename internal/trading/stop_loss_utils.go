// internal/trading/stop_loss_utils.go

package trading

import (
	"fmt"
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
func findExitDate(data []ml_stockdata.InMLDailyData, startDate time.Time, stopLossThreshold, trailingStopTriggerPrice, trailingStopTrigger, trailingStopUpdate float64) (time.Time, float64, error) {
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

		// ストップロスのチェックを先に行う
		if lowPrice <= stopLossThreshold || openPrice <= stopLossThreshold {
			endPrice = stopLossThreshold
			endDate = parsedDate
			break
		}

		// トレーリングストップのトリガーをチェック
		if closePrice >= trailingStopTriggerPrice {
			trailingStopTriggerPrice = round(closePrice * (1 + trailingStopTrigger/100))
			stopLossThreshold = round(closePrice * (1 - trailingStopUpdate/100))
		}
	}

	if endDate.IsZero() {
		// 最後のデータまで到達しても条件を満たさない場合、最終データを採用
		endPrice = data[len(data)-1].Close
		endDate, _ = parseDate(data[len(data)-1].Date)
	}
	fmt.Print(endPrice, "  ")
	return endDate, endPrice, nil
}

// round 関数: 四捨五入
func round(value float64) float64 {
	return float64(int(value*10+0.5)) / 10
}
