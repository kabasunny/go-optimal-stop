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
	bestTop3 := results[:3]
	worstTop3 := results[len(results)-3:]

	// 共通情報を表示
	if opts.ModelName != "" {
		fmt.Printf("モデル: %s,", opts.ModelName)
	}
	if opts.SignalCount > 0 {
		fmt.Printf("シグナル数: %d,", opts.SignalCount)
	}
	fmt.Printf("実行時間: %v\n", elapsedTime)

	// 結果を表示
	fmt.Println("  BEST 3:")
	for _, result := range bestTop3 {
		fmt.Printf("    [ LC(SL): %.2f%%, TST: %.2f%%, TS(SL): %.2f%%, 損益率: %.2f%%, 勝率: %.2f%% ]\n",
			result.StopLossPercentage, result.TrailingStopTrigger, result.TrailingStopUpdate, result.ProfitLoss, result.WinRate)
	}
	fmt.Println("  WORST 3:")
	for _, result := range worstTop3 {
		fmt.Printf("    [ LC(SL): %.2f%%, TST: %.2f%%, TS(SL): %.2f%%, 損益率: %.2f%%, 勝率: %.2f%% ]\n",
			result.StopLossPercentage, result.TrailingStopTrigger, result.TrailingStopUpdate, result.ProfitLoss, result.WinRate)
	}
}

// オプションを設定するための構造体と関数を定義
type resultOptions struct {
	ModelName   string
	SignalCount int
}

type ResultOption func(*resultOptions)

func WithModelName(name string) ResultOption {
	return func(opts *resultOptions) {
		opts.ModelName = name
	}
}

func WithSignalCount(count int) ResultOption {
	return func(opts *resultOptions) {
		opts.SignalCount = count
	}
}
