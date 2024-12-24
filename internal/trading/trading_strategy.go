// internal/trading/trading_strategy.go

package trading

import (
	"errors"
	"time"

	"go-optimal-stop/internal/stockdata"
)

// TradingStrategy 関数
func TradingStrategy(data *[]stockdata.Data, startDate time.Time, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (time.Time, time.Time, float64, error) {
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
