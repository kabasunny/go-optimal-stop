// internal/optimization/optimize_parameters.go

package optimization

import (
	"sort"
	"time"

	"go-optimal-stop/internal/stockdata"
	"go-optimal-stop/internal/trading"
)

// OptimizeParameters 関数の定義
func OptimizeParameters(inputData *[]stockdata.Data, tradeStartDate string) (stockdata.Result, stockdata.Result, []stockdata.Result) {
	var results []stockdata.Result

	stopLossPercentages := []float64{2.0, 3.0}
	trailingStopTriggers := []float64{5.0, 6.0, 7.0, 8.0, 9.0}
	trailingStopUpdates := []float64{2.0, 3.0}

	// tradeStartDate を time.Time 型に変換
	startDate, err := time.Parse("2006-01-02", tradeStartDate)
	if err != nil {
		return stockdata.Result{}, stockdata.Result{}, nil // エラーハンドリング
	}

	for _, stopLossPercentage := range stopLossPercentages {
		for _, trailingStopTrigger := range trailingStopTriggers {
			for _, trailingStopUpdate := range trailingStopUpdates {
				purchaseDate, exitDate, profitLoss, err := trading.TradingStrategy(inputData, startDate, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
				if err != nil {
					continue
				}
				result := stockdata.Result{
					StopLossPercentage:  stopLossPercentage,
					TrailingStopTrigger: trailingStopTrigger,
					TrailingStopUpdate:  trailingStopUpdate,
					ProfitLoss:          profitLoss,
					PurchaseDate:        purchaseDate.Format("2006-01-02"), // 文字列に変換
					ExitDate:            exitDate.Format("2006-01-02"),     // 文字列に変換
				}
				results = append(results, result)
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].ProfitLoss == results[j].ProfitLoss {
			if results[i].StopLossPercentage == results[j].StopLossPercentage {
				if results[i].TrailingStopTrigger == results[j].TrailingStopTrigger {
					return results[i].TrailingStopUpdate < results[j].TrailingStopUpdate
				}
				return results[i].TrailingStopTrigger < results[j].TrailingStopTrigger
			}
			return results[i].StopLossPercentage < results[j].StopLossPercentage
		}
		return results[i].ProfitLoss > results[j].ProfitLoss
	})

	bestResult := results[0]
	worstResult := results[len(results)-1]

	return bestResult, worstResult, results
}
