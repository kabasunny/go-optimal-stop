package trading

import (
	"errors"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
)

// TradingStrategy 関数
func TradingStrategy(response *ml_stockdata.InMLStockResponse, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (float64, float64, float64, float64, error) {
	totalProfitLoss := 0.0
	winCount := 0
	totalCount := 0
	currentPositiveStreak := 0.0
	maxPositiveStreak := 0.0
	currentNegativeStreak := 0.0
	maxNegativeStreak := 0.0

	for _, symbolData := range response.SymbolData {
		for _, signal := range symbolData.Signals {
			startDate, err := parseDate(signal)
			if err != nil {
				return 0, 0, 0, 0, err
			}

			// シンボルの株価データを使って最適化を実行
			_, _, profitLoss, _, _, err := singleTradingStrategy(&symbolData.DailyData, startDate, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
			if err != nil {
				return 0, 0, 0, 0, err
			}

			totalProfitLoss += profitLoss
			totalCount++

			if profitLoss > 0 {
				winCount++
				currentPositiveStreak += profitLoss
				if currentPositiveStreak > maxPositiveStreak {
					maxPositiveStreak = currentPositiveStreak
				}
				currentNegativeStreak = 0 // 負の連続をリセット
			} else {
				currentNegativeStreak += profitLoss
				if currentNegativeStreak < maxNegativeStreak {
					maxNegativeStreak = currentNegativeStreak
				}
				currentPositiveStreak = 0 // 正の連続をリセット
			}
		}
	}

	winRate := float64(winCount) / float64(totalCount) * 100
	return totalProfitLoss, winRate, maxPositiveStreak, maxNegativeStreak, nil
}

// singleTradingStrategy 関数
func singleTradingStrategy(data *[]ml_stockdata.InMLDailyData, startDate time.Time, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (time.Time, time.Time, float64, float64, float64, error) {
	d := *data
	if len(d) == 0 {
		return time.Time{}, time.Time{}, 0, 0, 0, errors.New("データが空です")
	}

	maxDate, err := parseDate(d[len(d)-1].Date)
	if err != nil {
		return time.Time{}, time.Time{}, 0, 0, 0, err
	}

	for !dateInData(d, startDate) {
		if startDate.After(maxDate) {
			return time.Time{}, time.Time{}, 0, 0, 0, errors.New("開始日がデータの範囲外です。無限ループを防ぐため、処理を中断")
		}
		startDate = startDate.AddDate(0, 0, 1)
	}

	purchaseDate, purchasePrice, err := findPurchaseDate(d, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, 0, 0, 0, err
	}

	stopLossThreshold, trailingStopTriggerPrice := calculateStopLoss(purchasePrice, stopLossPercentage, trailingStopTrigger)

	endDate, endPrice, err := findExitDate(d, startDate, stopLossThreshold, trailingStopTriggerPrice, trailingStopTrigger, trailingStopUpdate)
	if err != nil {
		return time.Time{}, time.Time{}, 0, 0, 0, err
	}

	profitLoss := (endPrice - purchasePrice) / purchasePrice * 100
	isProfit := profitLoss > 0
	profitLoss = round(profitLoss, isProfit)
	return purchaseDate, endDate, profitLoss, purchasePrice, endPrice, nil
}
