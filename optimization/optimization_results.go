package optimization

import (
	"fmt"
	"sort"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
)

func PrintResults(results []ml_stockdata.OptimizedResult, elapsedTime time.Duration, options ...ResultOption) {
	// オプションのデフォルト値を設定
	opts := &resultOptions{}
	for _, opt := range options {
		opt(opts)
	}

	// 結果をソート
	sort.Slice(results, func(i, j int) bool {
		return results[i].ProfitLoss > results[j].ProfitLoss
	})

	// トップ3の最良結果と最悪結果を取得
	bestTop3 := results[:5]
	worstTop3 := results[len(results)-5:]

	// 共通情報を表示
	if opts.ModelName != "" {
		fmt.Printf("モデル: %s , ", opts.ModelName)
	}
	if opts.SignalCount > 0 {
		fmt.Printf("シグナル数: %d , ", opts.SignalCount)
	}
	fmt.Printf("実行時間: %v\n", elapsedTime)

	// 結果を表示
	fmt.Println("  BEST 5:")
	for _, result := range bestTop3 {
		fmt.Printf("    [ LC: %.2f%%, TST: %.2f%%, TS: %.2f%%, 勝率: %.2f%%, 連続益: %.2f%%, 連続損: %.2f%%, 損益率: %.2f%% ]\n",
			result.StopLossPercentage, result.TrailingStopTrigger, result.TrailingStopUpdate, result.WinRate, result.MaxProfit, result.MaxLoss, result.ProfitLoss)
	}
	fmt.Println("  WORST 5:")
	for _, result := range worstTop3 {
		fmt.Printf("    [ LC: %.2f%%, TST: %.2f%%, TS: %.2f%%, 勝率: %.2f%%, 連続益: %.2f%%, 連続損: %.2f%%, 損益率: %.2f%% ]\n",
			result.StopLossPercentage, result.TrailingStopTrigger, result.TrailingStopUpdate, result.WinRate, result.MaxProfit, result.MaxLoss, result.ProfitLoss)
	}
}
