// main.go
package main

import (
	"flag"
	"fmt"
	"go-optimal-stop/experiment_proto"
	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/random_signals"
	"time"
)

func main() {
	// 引数を定義
	useRandom := flag.Bool("random", false, "Use random signals")
	flag.Parse()

	// Parameters構造体を作成し、関数を使ってパラメータを設定
	params := &ml_stockdata.Parameters{}
	params.SetStopLoss(2.0, 5.0, 1.0)
	params.SetTrailingStop(5.0, 10.0, 1.0)
	params.SetTrailingStopUpdate(2.0, 5.0, 1.0)

	if !*useRandom {
		fmt.Printf("学習モデルのシグナルで検証\n")

		filePath := "data/ml_stock_response/latest_response.bin"
		experiment_proto.RunOptimization(filePath, params)
	} else {
		fmt.Printf("ランダムにシグナルを作成し、結果を確認\n")

		csvDir := "../py-signal-buy/data/stock_data/formated_raw/2025-01-21"
		getSymbolsDir := "../py-signal-buy/data/stock_data/predictions/2025-01-21"
		symbols, err := random_signals.GetSymbolsFromCSVFiles(getSymbolsDir)
		if err != nil {
			fmt.Printf("Failed to get symbols from CSV files: %v\n", err)
			return
		}
		fmt.Printf("Symbols: %v\n", symbols)

		numSignals := 536
		// フラグを定義
		useRandomSeed := true // trueはランダム値、falseは固定値
		attempts := 5         // useRandomSeed := true の時、ランダム値試行を繰り返す回数

		// 本日の日付を取得し、365*2 さかのぼる
		today := time.Now()
		startDate := today.AddDate(-2, 0, 0).Format("2006-01-02")

		random_signals.RunRandomSignals(csvDir, symbols, numSignals, useRandomSeed, attempts, params, startDate)
	}
}

// 実行コマンド
// go run ./main.go --random # ランダムシグナルを使用
// go run ./main.go # 実践データを使用
