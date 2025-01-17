// main.go

package main

import (
	"flag"
	"go-optimal-stop/experiment_proto"
	"go-optimal-stop/random_signals"
)

func main() {
	// 引数を定義
	useRandom := flag.Bool("random", false, "Use random signals")
	flag.Parse()

	// 実践データを使用する場合
	if !*useRandom {
		filePath := "data/ml_stock_response/2025-01-17_11-12-14.bin"
		experiment_proto.RunOptimization(filePath)
	} else {
		// ランダムにシグナルを作成し、結果を確認
		csvDir := "../opti-ml-py/data/stock_data/formated_raw/2025-01-17"
		symbols := []string{"7203", "7267"}
		numSignals := 60 * len(symbols)
		seed := int64(42) // 固定シード
		random_signals.RunRandomSignals(csvDir, symbols, numSignals, seed)
	}
}

// 実行コマンド
// go run ./main.go --random # ランダムシグナルを使用
// go run ./main.go # 実践データを使用
