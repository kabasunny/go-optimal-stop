package optimization

import (
	"fmt"
	"sort"
	"sync"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/internal/trading"
)

// OptimizeParameters 関数は、与えられた株価データとトレーディングパラメータに基づいて最適なパラメータの組み合わせを見つける
func OptimizeParameters(response *ml_stockdata.InMLStockResponse, totalFunds *int, params *ml_stockdata.Parameters, commissionRate *float64) (ml_stockdata.OptimizedResult, ml_stockdata.OptimizedResult, []ml_stockdata.OptimizedResult) {
	// 結果を格納するスライスを事前に確保
	results := make([]ml_stockdata.OptimizedResult, 0,
		len(params.StopLossPercentages)*len(params.TrailingStopTriggers)*len(params.TrailingStopUpdates)*len(params.ATRMultipliers)*len(params.RiskPercentages))
	verbose := false

	var mu sync.Mutex     // 排他制御用のミューテックス
	var wg sync.WaitGroup // 同期用のWaitGroup

	// 各ストップロスパーセンテージ、トレーリングストップトリガー、トレーリングストップアップデート、ATR倍率、リスク許容度の組み合わせをループ処理
	for _, stopLossPercentage := range params.StopLossPercentages {
		for _, trailingStopTrigger := range params.TrailingStopTriggers {
			for _, trailingStopUpdate := range params.TrailingStopUpdates {
				for _, atrMultiplier := range params.ATRMultipliers {
					for _, riskPercentage := range params.RiskPercentages {
						// パラメータを設定
						param := ml_stockdata.Parameter{
							StopLossPercentage:  stopLossPercentage,
							TrailingStopTrigger: trailingStopTrigger,
							TrailingStopUpdate:  trailingStopUpdate,
							ATRMultiplier:       atrMultiplier,
							RiskPercentage:      riskPercentage,
						}

						wg.Add(1) // WaitGroupのカウントをインクリメント
						go func(param ml_stockdata.Parameter) {
							defer wg.Done() // 処理終了時にWaitGroupのカウントをデクリメント

							// トレーディング戦略を実行し、結果を取得
							result, err := trading.TradingStrategy(response, totalFunds, &param, commissionRate, verbose)
							if err != nil {
								fmt.Println("TradingStrategy error:", err) // エラーをログに出力
								mu.Lock()
								// エラー処理 (例: エラー数をカウント, 特定のエラーが発生したら処理を中断)
								mu.Unlock()
								return // エラーが発生した場合は、このゴルーチンを終了
							}

							// パラメータをOptimizedResult構造体に追加
							result.StopLossPercentage = param.StopLossPercentage
							result.TrailingStopTrigger = param.TrailingStopTrigger
							result.TrailingStopUpdate = param.TrailingStopUpdate
							result.ATRMultiplier = param.ATRMultiplier
							result.RiskPercentage = param.RiskPercentage

							// 結果をスライスに追加
							mu.Lock()                         // 排他制御開始
							results = append(results, result) // 結果をスライスに追加
							mu.Unlock()                       // 排他制御終了
						}(param)
					}
				}
			}
		}
	}

	wg.Wait() // すべてのGoルーチンの終了を待機

	// 結果をProfitLossでソート (最も重要な基準)
	sort.Slice(results, func(i, j int) bool {
		return results[i].ProfitLoss > results[j].ProfitLoss
	})

	// 最良の結果と最悪の結果を抽出
	bestResult := results[0]
	worstResult := results[len(results)-1]

	// 最良の結果、最悪の結果、全結果を返す
	return bestResult, worstResult, results
}
