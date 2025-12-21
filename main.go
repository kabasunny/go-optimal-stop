package main

import (
	"flag"
	"fmt"
	"go-optimal-stop/config" // Import the new config package
	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/optimization"
	"go-optimal-stop/random_signals"
	"log" // Import log for error handling
	"time"
)

func main() {
	start := time.Now() // 開始時刻を記録

	// Load configuration from config.yaml
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Define command-line arguments
	useRandom := flag.Bool("random", false, "Use random signals")
	dataFilePathFlag := flag.String("data", "", "Path to the .bin data file (overrides config)")
	outputFilePathFlag := flag.String("output", "", "Path to a JSON file to save optimization results")
	flag.Parse()

	// Determine the file path to use
	var actualFilePath string
	if *dataFilePathFlag != "" {
		actualFilePath = *dataFilePathFlag
	} else {
		actualFilePath = cfg.Simulation.DataFilePath
	}
	if actualFilePath == "" {
		log.Fatal("Data file path is not specified. Use --data flag or set data_file_path in config.yaml")
	}

	// Use values from config
	totalFunds := cfg.Simulation.TotalFunds
	commissionRate := cfg.Simulation.CommissionRate

	// Parameters構造体を作成し、関数を使ってパラメータを設定
	params := &ml_stockdata.Parameters{}
	params.SetStopLoss(
		cfg.OptimizationParameters.StopLoss.Min,
		cfg.OptimizationParameters.StopLoss.Max,
		cfg.OptimizationParameters.StopLoss.Step,
	)
	params.SetTrailingStop(
		cfg.OptimizationParameters.TrailingStop.Min,
		cfg.OptimizationParameters.TrailingStop.Max,
		cfg.OptimizationParameters.TrailingStop.Step,
	)
	params.SetTrailingStopUpdate(
		cfg.OptimizationParameters.TrailingStopUpdate.Min,
		cfg.OptimizationParameters.TrailingStopUpdate.Max,
		cfg.OptimizationParameters.TrailingStopUpdate.Step,
	)
	params.SetATRMultipliers(
		cfg.OptimizationParameters.ATRMultipliers.Min,
		cfg.OptimizationParameters.ATRMultipliers.Max,
		cfg.OptimizationParameters.ATRMultipliers.Step,
	)
	params.SetRiskPercentages(
		cfg.OptimizationParameters.RiskPercentages.Min,
		cfg.OptimizationParameters.RiskPercentages.Max,
		cfg.OptimizationParameters.RiskPercentages.Step,
	)

	// 総資金に対して、許容可能な最大ドローダウンを設定する

	if !*useRandom {
		fmt.Printf("学習モデルのシグナルで検証\n")

		// 許容ドローダウン値を渡す
		optimization.RunOptimization(&actualFilePath, &totalFunds, params, &commissionRate, outputFilePathFlag)
	} else {
		fmt.Printf("ランダムにシグナルを作成し、結果を確認\n")

		useRandomSeed := true // trueはランダム値、falseは固定値
		attempts := 1         // useRandomSeed := true の時、ランダム値試行を繰り返す回数

		// 許容ドローダウン値を渡す
		random_signals.RunRandomSignals(&actualFilePath, &totalFunds, params, &commissionRate, useRandomSeed, attempts, outputFilePathFlag)
	}

	elapsed := time.Since(start)         // 経過時間を計算
	fmt.Printf("全体の処理時間: %s\n", elapsed) // 経過時間を表示
}

// 実行コマンド
// go run ./main.go --random # ランダムシグナルを使用
// go run ./main.go # 実践データを使用
