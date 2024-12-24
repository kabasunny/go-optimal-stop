// internal/optimization/optimize_parameters.go

package optimization

import (
	"sort"

	"go-optimal-stop/internal/data"
	"go-optimal-stop/internal/trading"
)

// OptimizeParameters 関数の定義
func OptimizeParameters(data *[]data.Data, tradeStartDate string) (data.Result, data.Result, []data.Result) {
	var results []data.Result

	stopLossPercentages := []float64{2.0, 3.0}
	trailingStopTriggers := []float64{5.0, 6.0, 7.0, 8.0, 9.0}
	trailingStopUpdates := []float64{2.0, 3.0}

	for _, stopLossPercentage := range stopLossPercentages {
		for _, trailingStopTrigger := range trailingStopTriggers {
			for _, trailingStopUpdate := range trailingStopUpdates {
				purchaseDate, exitDate, profitLoss, err := trading.TradingStrategy(data, tradeStartDate, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
				if err != nil {
					continue
				}
				result := data.Result{
					StopLossPercentage:  stopLossPercentage,
					TrailingStopTrigger: trailingStopTrigger,
					TrailingStopUpdate:  trailingStopUpdate,
					ProfitLoss:          profitLoss,
					PurchaseDate:        purchaseDate,
					ExitDate:            exitDate,
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
