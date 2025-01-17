// random_signals.go

package random_signals

import (
	"fmt"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/optimization"
)

// 実行コマンド
// go run ./random_signals ランダムシード42を設定
// go run ./random_signals --random 完全にランダムにしたいとき
// go run ./random_signals/main.go ではmain.go ファイルのみをコンパイルして実行しようとするため動かない
func RunRandomSignals(csvDir string, symbols []string, numSignals int, seed int64) {
	startTime := time.Now() // 実行時間の測定開始

	// フラグを定義
	useRandomSeed := false // デフォルト値を設定

	var stockResponse ml_stockdata.InMLStockResponse
	var err error

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

	// Parameters構造体を作成し、関数を使ってパラメータを設定
	params := ml_stockdata.Parameters{}
	params.SetStopLoss(2.0, 5.0, 1.0)
	params.SetTrailingStop(5.0, 10.0, 1.0)
	params.SetTrailingStopUpdate(2.0, 5.0, 1.0)

	// 総試行回数を算出
	totalTrials := len(params.StopLossPercentages) * len(params.TrailingStopTriggers) * len(params.TrailingStopUpdates) * len(stockResponse.SymbolData) * numSignals
	fmt.Printf("総試行回数: %d\n", totalTrials)

	// パラメータの最適化を実行
	_, _, results := optimization.OptimizeParameters(&stockResponse, params)

	// 実行時間を測定
	elapsedTime := time.Since(startTime)

	// 結果を表示
	optimization.PrintOverallResults(results, elapsedTime)
}
