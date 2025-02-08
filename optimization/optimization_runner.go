package optimization

import (
	"fmt"
	"os"
	"time"

	"go-optimal-stop/experiment_proto"
	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/internal/trading"

	"google.golang.org/protobuf/proto"
)

// 許容ドローダウン値を渡す
func RunOptimization(filePath *string, totalFunds *int, params *ml_stockdata.Parameters, commissionRate *float64) {
	// startTime := time.Now() // 実行時間の測定開始

	// ファイルを読み込み、stockResponseにプロトコルバッファバイナリからデータをマッピング
	data, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("ファイルの読み込みエラー: %v\n", err)
		return
	}

	var protoResponse experiment_proto.MLStockResponse
	if err := proto.Unmarshal(data, &protoResponse); err != nil {
		fmt.Printf("プロトコルバッファのアンマーシャルエラー: %v\n", err)
		return
	}

	// プロトコルバッファから内部MLStockResponse型への変換
	stockResponse := experiment_proto.ConvertProtoToInternal(&protoResponse)

	// protoResponse 内のシンボルのリストを表示
	var symbols []string
	for _, symbolData := range protoResponse.GetSymbolData() {
		symbols = append(symbols, symbolData.Symbol)
	}
	fmt.Printf("Symbols: %v\n", symbols)

	// protoResponse 内の全シンボルの全シグナル数の合計を取得
	numSignals := 0
	for _, symbolData := range stockResponse.SymbolData {
		numSignals += len(symbolData.Signals)
	}

	// 総試行回数を算出
	trials := len(params.StopLossPercentages) * len(params.TrailingStopTriggers) * len(params.TrailingStopUpdates) * len(params.ATRMultipliers) * len(params.RiskPercentages)
	totalTrials := trials * numSignals * len(stockResponse.SymbolData)
	fmt.Printf("パラメタ組合せ: %d, シグナル数: %d, 総試行回数: %d\n", trials, numSignals, totalTrials)

	// モデル名correct_labelを最初に配置し、あとはmodel_predictionsフィールドから自動抽出し
	modelNames := []string{"correct_label"}
	otherModelNames := []string{}
	for modelName := range protoResponse.SymbolData[0].ModelPredictions {
		if modelName != "correct_label" && modelName != "ensemble_label" {
			otherModelNames = append(otherModelNames, modelName)
		}
	}
	modelNames = append(modelNames, otherModelNames...)
	modelNames = append(modelNames, "ensemble_label")
	fmt.Printf("実行SIM一覧: %v\n", modelNames)

	for _, modelName := range modelNames {
		// すべてのシンボルに対してシグナルを一斉に置き換える
		originalSignals := make(map[int][]string)
		var signalCount int
		for i := range protoResponse.SymbolData {
			if modelPredictions, ok := protoResponse.SymbolData[i].ModelPredictions[modelName]; ok && modelPredictions != nil {
				originalSignals[i] = stockResponse.SymbolData[i].Signals
				stockResponse.SymbolData[i].Signals = modelPredictions.PredictionDates // 新しいシグナルを設定
				signalCount += len(modelPredictions.PredictionDates)                   // 各シンボルのシグナル数をカウント
			} else {
				fmt.Printf("モデル: %s の予測データが見つかりませんでした。シンボル: %s をスキップします。\n", modelName, protoResponse.SymbolData[i].Symbol)
			}
		}

		// モデルの最適化開始時間を記録
		modelStartTime := time.Now()

		// すべてのシグナルが置き換えられた後にパラメータの最適化を実行
		_, _, modelResults := OptimizeParameters(&stockResponse, totalFunds, params, commissionRate)

		// モデルの実行時間を測定
		modelElapsedTime := time.Since(modelStartTime)

		// モデルごとの結果を表示
		bestparm, worstparam, _ := PrintAndReturnResults(modelResults, modelElapsedTime, WithModelName(modelName), WithSignalCount(signalCount))

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

		// 元のシグナルに戻す
		for i := range protoResponse.SymbolData {
			stockResponse.SymbolData[i].Signals = originalSignals[i]
		}
	}
}
