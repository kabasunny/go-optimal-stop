// cmd/main.go

package main

import (
	"flag"
	"fmt"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/internal/optimization"
)

// 実行コマンド
// go run ./cmd ランダムシード42を設定
// go run ./cmd --random 完全にランダムにしたいとき
// go run ./cmd/main.go ではmain.go ファイルのみをコンパイルして実行しようとするため動かない
func main() {
	startTime := time.Now() // 実行時間の測定開始

	// フラグを定義
	useRandomSeed := flag.Bool("random", false, "Use random seed")
	flag.Parse()

	csvDir := "csv" // CSVファイルが保存されているディレクトリ
	symbols := []string{
		"7203.T", // Toyota Motor Corporation
		"7201.T", // Nissan Motor Co., Ltd.
		"7267.T", // Honda Motor Co., Ltd.
		"7261.T", // Mazda Motor Corporation
		"7269.T", // Suzuki Motor Corporation
		// "7262.T",  // Mitsubishi Motors Corporation 上場廃止
		"7270.T", // Subaru Corporation
		"7202.T", // Isuzu Motors Limited
		"7205.T", // Hino Motors, Ltd.
		"7211.T", // Mitsubishi Fuso Truck and Bus Corporation
		"7224.T", // Shizuoka Daihatsu Motor Co., Ltd.
		"7266.T", // Showa Corporation} // シンボルのリスト
	}
	numSignals := 50 // ランダムに選ぶシグナルの数

	var stockResponse ml_stockdata.MLStockResponse
	var err error

	if *useRandomSeed {
		// 完全にランダムにシグナルを生成
		stockResponse, err = createStockResponse(csvDir, symbols, numSignals)
	} else {
		// 固定シードを使用してシグナルを生成
		seed := int64(42)
		stockResponse, err = createStockResponse(csvDir, symbols, numSignals, seed)
	}

	if err != nil {
		fmt.Printf("StockResponseの作成エラー: %v\n", err)
		return
	}

	// Parameters構造体を作成し、関数を使ってパラメータを設定
	params := ml_stockdata.Parameters{}
	params.SetStopLoss(2.0, 10.0, 1.0)
	params.SetTrailingStop(5.0, 20.0, 1.0)
	params.SetTrailingStopUpdate(2.0, 10.0, 1.0)

	// 総試行回数を算出
	totalTrials := len(params.StopLossPercentages) * len(params.TrailingStopTriggers) * len(params.TrailingStopUpdates) * len(stockResponse.SymbolData) * numSignals
	fmt.Printf("総試行回数: %d\n", totalTrials)

	// パラメータの最適化を実行
	bestResult, worstResult, _ := optimization.OptimizeParameters(&stockResponse, params)

	// 実行時間を測定
	elapsedTime := time.Since(startTime)

	// 結果を表示
	fmt.Printf("最良の結果: %+v\n", bestResult)
	fmt.Printf("最悪の結果: %+v\n", worstResult)
	fmt.Printf("実行時間: %v\n", elapsedTime)
}
