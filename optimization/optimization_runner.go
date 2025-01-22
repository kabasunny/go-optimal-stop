package optimization

import (
	"fmt"
	"os"
	"time"

	"go-optimal-stop/experiment_proto"
	"go-optimal-stop/internal/ml_stockdata"

	"google.golang.org/protobuf/proto"
)

func RunOptimization(filePath string, params *ml_stockdata.Parameters) {
	startTime := time.Now() // 実行時間の測定開始

	// ファイルを読み込み、stockResponseにプロトコルバッファバイナリからデータをマッピング
	data, err := os.ReadFile(filePath)
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
	trials := len(params.StopLossPercentages) * len(params.TrailingStopTriggers) * len(params.TrailingStopUpdates) * len(stockResponse.SymbolData)
	totalTrials := trials * numSignals
	fmt.Printf("パラメタ組合せ: %d, 正解ラベル数: %d, 総試行回数: %d\n", trials, numSignals, totalTrials)

	// パラメータの最適化を実行
	_, _, results := OptimizeParameters(&stockResponse, params)

	// 実行時間を測定
	elapsedTime := time.Since(startTime)

	// 結果を表示
	PrintResults(results, elapsedTime)

	// モデル名をmodel_predictionsフィールドから自動抽出
	modelNames := extractModelNames(&protoResponse)
	fmt.Printf("シミュレーションモデル名: %v\n", modelNames)

	for _, modelName := range modelNames {
		// モデルの予測データを取得
		if modelPredictions, ok := protoResponse.SymbolData[0].ModelPredictions[modelName]; ok && modelPredictions != nil {
			// モデルごとのシグナルを設定
			modelSignals := modelPredictions.PredictionDates
			// シグナルを設定する前に元のシグナルを保存しておく（他のモデルで使うため）
			originalSignals := stockResponse.SymbolData[0].Signals
			stockResponse.SymbolData[0].Signals = modelSignals

			// シグナル数を取得
			modelSignalCount := len(modelSignals)

			// モデルの最適化開始時間を記録
			modelStartTime := time.Now()

			// パラメータの最適化を再実行
			_, _, modelResults := OptimizeParameters(&stockResponse, params)

			// モデルの実行時間を測定
			modelElapsedTime := time.Since(modelStartTime)

			// モデルごとの結果を表示
			PrintResults(modelResults, modelElapsedTime, WithModelName(modelName), WithSignalCount(modelSignalCount))

			// 元のシグナルに戻す
			stockResponse.SymbolData[0].Signals = originalSignals
		} else {
			fmt.Printf("モデル: %s の予測データが見つかりませんでした。スキップします。\n", modelName)
		}
	}
}
