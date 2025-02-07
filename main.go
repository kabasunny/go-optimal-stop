package main

import (
	"flag"
	"fmt"
	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/optimization"
	"go-optimal-stop/random_signals"
	"time"
)

func main() {
	start := time.Now() // 開始時刻を記録

	filePath := "../py-signal-buy/result/ml_stock_response/proto_kmeans-cluster_label_4.bin"

	totalFunds := 2000000
	commissionRate := 0.2 // 手数料率（例: 0.2%）

	// 引数を定義
	useRandom := flag.Bool("random", false, "Use random signals")
	flag.Parse()

	// Parameters構造体を作成し、関数を使ってパラメータを設定
	params := &ml_stockdata.Parameters{}
	params.SetStopLoss(1.0, 5.0, 1.0)
	params.SetTrailingStop(5.0, 20.0, 1.0)
	params.SetTrailingStopUpdate(1.0, 5.0, 1.0)

	// 総資金に対して、許容可能な最大ドローダウンを設定する

	if !*useRandom {
		fmt.Printf("学習モデルのシグナルで検証\n")

		// 許容ドローダウン値を渡す
		optimization.RunOptimization(&filePath, &totalFunds, params, &commissionRate)
	} else {
		fmt.Printf("ランダムにシグナルを作成し、結果を確認\n")

		useRandomSeed := true // trueはランダム値、falseは固定値
		attempts := 1         // useRandomSeed := true の時、ランダム値試行を繰り返す回数

		// 許容ドローダウン値を渡す
		random_signals.RunRandomSignals(&filePath, &totalFunds, params, &commissionRate, useRandomSeed, attempts)
	}

	elapsed := time.Since(start)         // 経過時間を計算
	fmt.Printf("全体の処理時間: %s\n", elapsed) // 経過時間を表示
}

// 実行コマンド
// go run ./main.go --random # ランダムシグナルを使用
// go run ./main.go # 実践データを使用
