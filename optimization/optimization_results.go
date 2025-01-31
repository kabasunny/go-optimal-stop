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

	// トップ5の最良結果と最悪結果を取得
	bestTop5 := results[:5]
	worstTop5 := results[len(results)-5:]

	// 共通情報を表示
	if opts.ModelName != "" {
		fmt.Printf("実施SIM名: %s, ", opts.ModelName)
	}
	if opts.SignalCount > 0 {
		fmt.Printf("シグナル数: %d, ", opts.SignalCount)
	}
	fmt.Printf("実行時間: %v\n", elapsedTime)

	// ラベルの説明を表示
	fmt.Println("LC:ロスカット, TT:トレリングストップトリガ, TU:TS更新値, WR:勝率, CP:連続益, CL:連続損, PL:損益率, TW:総勝数, TL:総負数, AP:平均益, AL:平均損, MD:最大ドローダウン, SR:シャープレシオ, RR:リスクリワード, EV:期待値")

	// 結果を表示
	fmt.Println("  BEST 5:")
	for _, result := range bestTop5 {
		fmt.Printf("    [ LC: %.1f%%, TT: %.1f%%, TU: %.1f%%, WR: %.1f%%, CP: %.1f%%, CL: %.1f%%, PL: %.1f%%, TW: %d, TL: %d, AP: %.1f%%, AL: %.1f%%, MD: %.1f%%, SR: %.1f, RR: %.1f, EV: %.1f%% ]\n",
			result.StopLossPercentage, result.TrailingStopTrigger, result.TrailingStopUpdate, result.WinRate, result.MaxConsecutiveProfit, result.MaxConsecutiveLoss, result.ProfitLoss,
			result.TotalWins, result.TotalLosses, result.AverageProfit, result.AverageLoss, result.MaxDrawdown, result.SharpeRatio, result.RiskRewardRatio, result.ExpectedValue)
	}
	fmt.Println("  WORST 5:")
	for _, result := range worstTop5 {
		fmt.Printf("    [ LC: %.1f%%, TT: %.1f%%, TU: %.1f%%, WR: %.1f%%, CP: %.1f%%, CL: %.1f%%, PL: %.1f%%, TW: %d, TL: %d, AP: %.1f%%, AL: %.1f%%, MD: %.1f%%, SR: %.1f, RR: %.1f, EV: %.1f%% ]\n",
			result.StopLossPercentage, result.TrailingStopTrigger, result.TrailingStopUpdate, result.WinRate, result.MaxConsecutiveProfit, result.MaxConsecutiveLoss, result.ProfitLoss,
			result.TotalWins, result.TotalLosses, result.AverageProfit, result.AverageLoss, result.MaxDrawdown, result.SharpeRatio, result.RiskRewardRatio, result.ExpectedValue)
	}
}
