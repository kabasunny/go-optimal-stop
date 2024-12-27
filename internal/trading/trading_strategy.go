// internal/trading/trading_strategy.go

package trading

import (
	"errors"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
)

// TradingStrategy 関数
func TradingStrategy(response *ml_stockdata.MLStockResponse, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (float64, error) {
	totalProfitLoss := 0.0

	for _, symbolData := range response.SymbolData {
		for _, signal := range symbolData.Signals {
			startDate, err := parseDate(signal)
			if err != nil {
				return 0, err
			}

			// シンボルの株価データを使って最適化を実行
			_, _, profitLoss, err := singleTradingStrategy(&symbolData.DailyData, startDate, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
			if err != nil {
				return 0, err
			}
			totalProfitLoss += profitLoss
		}
	}

	return totalProfitLoss, nil
}

// singleTradingStrategy 関数
func singleTradingStrategy(data *[]ml_stockdata.MLDailyData, startDate time.Time, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (time.Time, time.Time, float64, error) {
	d := *data
	if len(d) == 0 {
		return time.Time{}, time.Time{}, 0, errors.New("データが空です")
	}

	maxDate, err := parseDate(d[len(d)-1].Date)
	if err != nil {
		return time.Time{}, time.Time{}, 0, err
	}

	for !dateInData(d, startDate) {
		if startDate.After(maxDate) {
			return time.Time{}, time.Time{}, 0, errors.New("開始日がデータの範囲外です。無限ループを防ぐため、処理を中断")
		}
		startDate = startDate.AddDate(0, 0, 1)
	}

	purchaseDate, purchasePrice, err := findPurchaseDate(d, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, 0, err
	}

	stopLossThreshold, trailingStopTriggerPrice := calculateStopLoss(purchasePrice, stopLossPercentage, trailingStopTrigger)

	endDate, endPrice, err := findExitDate(d, startDate, stopLossThreshold, trailingStopTriggerPrice, trailingStopTrigger, trailingStopUpdate)
	if err != nil {
		return time.Time{}, time.Time{}, 0, err
	}

	profitLoss := round((endPrice - purchasePrice) / purchasePrice * 100)
	return purchaseDate, endDate, profitLoss, nil
}
