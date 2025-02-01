package trading

import (
	"go-optimal-stop/internal/ml_stockdata"
	"time"
)

// roundUp 関数: 四捨五入（切り上げ）
func roundUp(value float64) float64 {
	return float64(int(value*10+1)) / 10
}

// roundDown 関数: 四捨五入（切り捨て）
func roundDown(value float64) float64 {
	return float64(int(value*10)) / 10
}

// findExitDate: 売却日と売却価格を決定
func findExitDate(data []ml_stockdata.InMLDailyData, purchaseDate time.Time, purchasePrice, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (time.Time, float64, error) {
	var endDate time.Time
	var endPrice float64

	// ストップロスとトレーリングストップの閾値を計算
	stopLossThreshold := roundDown(purchasePrice * (1 - stopLossPercentage/100))
	trailingStopTriggerPrice := roundUp(purchasePrice * (1 + trailingStopTrigger/100))

	// purchaseDate 以降のデータを取得（スライスを最適化）
	var startIndex int
	for i, day := range data {
		parsedDate, err := parseDate(day.Date)
		if err != nil {
			return time.Time{}, 0, err
		}
		if !parsedDate.Before(purchaseDate) {
			startIndex = i
			break
		}
	}
	filteredData := data[startIndex:]

	// トレーリングストップの監視
	for _, day := range filteredData {
		parsedDate, err := parseDate(day.Date)
		if err != nil {
			return time.Time{}, 0, err
		}

		openPrice := day.Open
		lowPrice := day.Low
		closePrice := day.Close

		// ストップロスに到達した場合
		if lowPrice <= stopLossThreshold || openPrice <= stopLossThreshold {
			endPrice = stopLossThreshold
			endDate = parsedDate
			break
		}

		// トレーリングストップのトリガーをチェック
		if closePrice >= trailingStopTriggerPrice {
			trailingStopTriggerPrice = roundUp(closePrice * (1 + trailingStopTrigger/100))
			stopLossThreshold = roundDown(closePrice * (1 - trailingStopUpdate/100))
		}
	}

	// 途中で売却しなかった場合、最終データを採用
	if endDate.IsZero() {
		lastIndex := len(filteredData) - 1
		endPrice = filteredData[lastIndex].Close
		endDate, _ = parseDate(filteredData[lastIndex].Date)
	}

	return endDate, endPrice, nil
}
