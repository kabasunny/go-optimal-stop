package optimization

import (
	"sort"
	"sync"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/internal/trading"
)

// OptimizeParameters 関数は、与えられた株価データとトレーディングパラメータに基づいて最適なパラメータの組み合わせを見つける
func OptimizeParameters(response *ml_stockdata.InMLStockResponse, totalFunds *int, params *ml_stockdata.Parameters) (ml_stockdata.OptimizedResult, ml_stockdata.OptimizedResult, []ml_stockdata.OptimizedResult) {
	var results []ml_stockdata.OptimizedResult // 最適化結果を格納するスライス
	var mu sync.Mutex                          // 排他制御用のミューテックス
	var wg sync.WaitGroup                      // 同期用のWaitGroup

	// 各ストップロスパーセンテージ、トレーリングストップトリガー、トレーリングストップアップデートの組み合わせをループ処理
	for _, stopLossPercentage := range params.StopLossPercentages {
		for _, trailingStopTrigger := range params.TrailingStopTriggers {
			for _, trailingStopUpdate := range params.TrailingStopUpdates {
				wg.Add(1) // WaitGroupのカウントをインクリメント
				go func(totalFunds *int, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) {
					defer wg.Done() // 処理終了時にWaitGroupのカウントをデクリメント
					// トレーディング戦略を実行し、結果を取得
					result, err := trading.TradingStrategy(response, totalFunds, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
					if err != nil {
						return
					}
					// パラメータをOptimizedResult構造体に追加
					result.StopLossPercentage = stopLossPercentage
					result.TrailingStopTrigger = trailingStopTrigger
					result.TrailingStopUpdate = trailingStopUpdate

					// 結果をスライスに追加
					mu.Lock()                         // 排他制御開始
					results = append(results, result) // 結果をスライスに追加
					mu.Unlock()                       // 排他制御終了
				}(totalFunds, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
			}
		}
	}

	wg.Wait() // すべてのGoルーチンの終了を待機

	// 結果をProfitLoss、StopLossPercentage、TrailingStopTrigger、TrailingStopUpdateでソート
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

	// 最良の結果と最悪の結果を抽出
	bestResult := results[0]
	worstResult := results[len(results)-1]

	// 最良の結果、最悪の結果、全結果を返す
	return bestResult, worstResult, results
}
