package optimization

import (
	"fmt"
	"sort"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
)

func PrintOverallResults(results []ml_stockdata.OptimizedionResult, elapsedTime time.Duration) {
	// 結果をソート
	sort.Slice(results, func(i, j int) bool {
		return results[i].ProfitLoss > results[j].ProfitLoss
	})

	// トップ3の最良結果と最悪結果を取得
	bestTop3 := results[:3]
	worstTop3 := results[len(results)-3:]

	// 結果を表示
	fmt.Printf("全体の結果:\n")
	fmt.Println("BEST 3:")
	for _, result := range bestTop3 {
		fmt.Printf("  [ LC(SL): %.2f%%, TST: %.2f%%, TS(SL): %.2f%%, ProfitLoss: %.2f%% ]\n",
			result.StopLossPercentage, result.TrailingStopTrigger, result.TrailingStopUpdate, result.ProfitLoss)
	}
	fmt.Println("WORST 3:")
	for _, result := range worstTop3 {
		fmt.Printf("  [ LC(SL): %.2f%%, TST: %.2f%%, TS(SL): %.2f%%, ProfitLoss: %.2f%% ]\n",
			result.StopLossPercentage, result.TrailingStopTrigger, result.TrailingStopUpdate, result.ProfitLoss)
	}
	fmt.Printf("実行時間: %v\n", elapsedTime)
}

func PrintModelResults(modelName string, modelSignalCount int, results []ml_stockdata.OptimizedionResult, modelElapsedTime time.Duration) {
	// 結果をソート
	sort.Slice(results, func(i, j int) bool {
		return results[i].ProfitLoss > results[j].ProfitLoss
	})

	// トップ3の最良結果と最悪結果を取得
	modelBestTop3 := results[:3]
	modelWorstTop3 := results[len(results)-3:]

	// モデルごとの結果を表示
	fmt.Printf("  シグナル数: %d, 実行時間: %v\n", modelSignalCount, modelElapsedTime)
	fmt.Println("  BEST 3:")
	for _, result := range modelBestTop3 {
		fmt.Printf("    [ LC(SL): %.2f%%, TST: %.2f%%, TS(SL): %.2f%%, ProfitLoss: %.2f%% ]\n",
			result.StopLossPercentage, result.TrailingStopTrigger, result.TrailingStopUpdate, result.ProfitLoss)
	}
	fmt.Println("  WORST 3:")
	for _, result := range modelWorstTop3 {
		fmt.Printf("    [ LC(SL): %.2f%%, TST: %.2f%%, TS(SL): %.2f%%, ProfitLoss: %.2f%% ]\n",
			result.StopLossPercentage, result.TrailingStopTrigger, result.TrailingStopUpdate, result.ProfitLoss)
	}
}
