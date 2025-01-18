// main.go

package main

import (
	"flag"
	"fmt"
	"go-optimal-stop/experiment_proto"
	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/random_signals"
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

		csvDir := "../opti-ml-py/data/stock_data/formated_raw/2025-01-17"
		// csvDir := "."
		symbols := []string{"7203", "7267"}
		numSignals := 447
		// フラグを定義
		useRandomSeed := true // trueはランダム値、falseは固定値
		attempts := 5         // useRandomSeed := true の時、ランダム値試行を繰り返す回数
		random_signals.RunRandomSignals(csvDir, symbols, numSignals, useRandomSeed, attempts, params)
	}
}

// 実行コマンド
// go run ./main.go --random # ランダムシグナルを使用
// go run ./main.go # 実践データを使用
