package random_signals

import (
	"fmt"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/optimization"
)

func RunRandomSignals(filePath *string, totalFunds *int, params *ml_stockdata.Parameters, useRandomSeed bool, attempts int) {

	var stockResponse ml_stockdata.InMLStockResponse
	var err error
	var numSignals int
	var symbols []string
	seed := int64(42) // 固定シード

	for i := 0; i < attempts; i++ {
		startTime := time.Now() // 実行時間の測定開始
		fmt.Printf("ランダム試行 %d 回目 / %d 回中\n", i+1, attempts)
		if useRandomSeed {
			// 完全にランダムにシグナルを生成
			stockResponse, numSignals, symbols, err = createStockResponse(filePath)

		} else {
			// 固定シードを使用してシグナルを生成
			stockResponse, numSignals, symbols, err = createStockResponse(filePath, seed)
		}

		if err != nil {
			fmt.Printf("StockResponseの作成エラー: %v\n", err)
			return
		}

		fmt.Printf("Symbols: %v\n", symbols)

		// 総試行回数を算出
		trials := len(params.StopLossPercentages) * len(params.TrailingStopTriggers) * len(params.TrailingStopUpdates)
		totalTrials := trials * numSignals * len(stockResponse.SymbolData)
		fmt.Printf("パラメタ組合せ: %d, シグナル数: %d, 総試行回数: %d\n", trials, numSignals, totalTrials)

		// パラメータの最適化を実行
		_, _, results := optimization.OptimizeParameters(&stockResponse, totalFunds, params)

		// 実行時間を測定
		elapsedTime := time.Since(startTime)

		// 結果を表示
		_, _ = optimization.PrintAndReturnResults(results, elapsedTime)
	}
}
