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

	// 引数を定義
	useRandom := flag.Bool("random", false, "Use random signals")
	flag.Parse()

	// Parameters構造体を作成し、関数を使ってパラメータを設定
	params := &ml_stockdata.Parameters{}
	params.SetStopLoss(2.0, 8.0, 1.0)
	params.SetTrailingStop(5.0, 15.0, 1.0)
	params.SetTrailingStopUpdate(2.0, 8.0, 1.0)

	if !*useRandom {
		fmt.Printf("学習モデルのシグナルで検証\n")

		filePath := "data/ml_stock_response/proto_kmeans-cluster_label_0.bin"
		optimization.RunOptimization(filePath, params)
	} else {
		fmt.Printf("ランダムにシグナルを作成し、結果を確認\n")

		// 今日の日付を取得し、フォーマットする
		today := time.Now().Format("2006-01-02")

		// 予めPython側が、今日の日付で出力している前提
		csvDir := fmt.Sprintf("../py-signal-buy/data/stock_data/formated_raw/%s", today)
		getSymbolsDir := fmt.Sprintf("../py-signal-buy/data/stock_data/predictions/%s", today)

		symbols, err := random_signals.GetSymbolsFromCSVFiles(getSymbolsDir)
		if err != nil {
			fmt.Printf("Failed to get symbols from CSV files: %v\n", err)
			return
		}
		fmt.Printf("Symbols: %v\n", symbols)

		numSignals := 16513
		// フラグを定義
		useRandomSeed := true // trueはランダム値、falseは固定値
		attempts := 3         // useRandomSeed := true の時、ランダム値試行を繰り返す回数

		// 本日の日付を取得し、365*2 さかのぼる
		startDate := time.Now().AddDate(-2, 0, 0).Format("2006-01-02")

		random_signals.RunRandomSignals(csvDir, symbols, numSignals, useRandomSeed, attempts, params, startDate)
	}

	elapsed := time.Since(start)         // 経過時間を計算
	fmt.Printf("全体の処理時間: %s\n", elapsed) // 経過時間を表示
}

// 実行コマンド
// go run ./main.go --random # ランダムシグナルを使用
// go run ./main.go # 実践データを使用
