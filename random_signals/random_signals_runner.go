// random_signals.go

package random_signals

import (
	"fmt"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/optimization"
)

func RunRandomSignals(csvDir string, symbols []string, numSignals int, useRandomSeed bool, attempts int, params *ml_stockdata.Parameters) {

	var stockResponse ml_stockdata.InMLStockResponse
	var err error
	seed := int64(42) // 固定シード

	for i := 0; i < attempts; i++ {
		startTime := time.Now() // 実行時間の測定開始
		fmt.Printf("ランダム試行 %d 回目 / %d 回中\n", i, attempts)
		if useRandomSeed {
			// 完全にランダムにシグナルを生成
			stockResponse, err = createStockResponse(csvDir, symbols, numSignals)

		} else {
			// 固定シードを使用してシグナルを生成
			stockResponse, err = createStockResponse(csvDir, symbols, numSignals, seed)
		}

		if err != nil {
			fmt.Printf("StockResponseの作成エラー: %v\n", err)
			return
		}

		// 総試行回数を算出
		totalTrials := len(params.StopLossPercentages) * len(params.TrailingStopTriggers) * len(params.TrailingStopUpdates) * len(stockResponse.SymbolData) * numSignals
		fmt.Printf("総試行回数: %d, シグナル数: %d\n", totalTrials, numSignals)

		// パラメータの最適化を実行
		_, _, results := optimization.OptimizeParameters(&stockResponse, params)

		// 実行時間を測定
		elapsedTime := time.Since(startTime)

		// 結果を表示
		optimization.PrintResults(results, elapsedTime)
	}
}
