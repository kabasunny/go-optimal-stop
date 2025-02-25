package random_signals

import (
	"fmt"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/internal/trading"
	"go-optimal-stop/optimization"
)

func RunRandomSignals(filePath *string, totalFunds *int, params *ml_stockdata.Parameters, commissionRate *float64, useRandomSeed bool, attempts int) {

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
		trials := len(params.StopLossPercentages) * len(params.TrailingStopTriggers) * len(params.TrailingStopUpdates) * len(params.ATRMultipliers) * len(params.RiskPercentages)
		totalTrials := trials * numSignals * len(stockResponse.SymbolData)
		fmt.Printf("パラメタ組合せ: %d, シグナル数: %d, 総試行回数: %d\n", trials, numSignals, totalTrials)

		// パラメータの最適化を実行
		_, _, results := optimization.OptimizeParameters(&stockResponse, totalFunds, params, commissionRate)

		// 実行時間を測定
		elapsedTime := time.Since(startTime)

		// 結果を表示
		bestparm, worstparam, _ := optimization.PrintAndReturnResults(results, elapsedTime)

		verbose := true
		if verbose {
			fmt.Println("BESTパラメータで、トレードシミュレーション")
			fmt.Printf(" [%-2s](%9s) %9s : %7s - %7s (%5s)[ %9s (%4s) - %9s ] %6s, %6s, %6s\n",
				"銘柄", "entry日", "exit日", "entry株価", "exit株価", "size", "entry金額", "総割合", "exit金額", "単損益", "総損益", "総資金")
			_, _ = trading.TradingStrategy(&stockResponse, totalFunds, &bestparm, commissionRate, verbose)

			fmt.Println("WORSTパラメータで、トレードシミュレーション")
			fmt.Printf(" [%-2s](%9s) %9s : %7s - %7s (%5s)[ %9s (%4s) - %9s ] %6s, %6s, %6s\n",
				"銘柄", "entry日", "exit日", "entry株価", "exit株価", "size", "entry金額", "総割合", "exit金額", "単損益", "総損益", "総資金")
			_, _ = trading.TradingStrategy(&stockResponse, totalFunds, &worstparam, commissionRate, verbose)
		}

	}
}
