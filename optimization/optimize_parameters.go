// internal/optimization/optimize_parameters.go

package optimization

import (
	"sort"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/internal/trading"
)

// OptimizeParameters 関数の定義
func OptimizeParameters(response *ml_stockdata.InMLStockResponse, params *ml_stockdata.Parameters) (ml_stockdata.OptimizedResult, ml_stockdata.OptimizedResult, []ml_stockdata.OptimizedResult) {
	var results []ml_stockdata.OptimizedResult

	for _, stopLossPercentage := range params.StopLossPercentages {
		for _, trailingStopTrigger := range params.TrailingStopTriggers {
			for _, trailingStopUpdate := range params.TrailingStopUpdates {
				// TradingStrategy 関数を呼び出して総利益を計算
				totalProfitLoss, winRate, maxProfit, maxLoss, err := trading.TradingStrategy(response, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
				if err != nil {
					continue
				}
				result := ml_stockdata.OptimizedResult{
					StopLossPercentage:  stopLossPercentage,
					TrailingStopTrigger: trailingStopTrigger,
					TrailingStopUpdate:  trailingStopUpdate,
					ProfitLoss:          totalProfitLoss,
					WinRate:             winRate,
					MaxProfit:           maxProfit,
					MaxLoss:             maxLoss,
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
