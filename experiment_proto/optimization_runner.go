package experiment_proto

import (
	"fmt"
	"os"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/optimization"

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

	var protoResponse MLStockResponse
	if err := proto.Unmarshal(data, &protoResponse); err != nil {
		fmt.Printf("プロトコルバッファのアンマーシャルエラー: %v\n", err)
		return
	}

	// プロトコルバッファから内部MLStockResponse型への変換
	stockResponse := ConvertProtoToInternal(&protoResponse)

	// protoResponse 内の全シンボルの全シグナル数の合計を取得
	numSignals := 0
	for _, symbolData := range stockResponse.SymbolData {
		numSignals += len(symbolData.Signals)
	}

	// 総試行回数を算出
	totalTrials := len(params.StopLossPercentages) * len(params.TrailingStopTriggers) * len(params.TrailingStopUpdates) * len(stockResponse.SymbolData) * numSignals
	fmt.Printf("サーチ回数/試行: %d, シグナル数: %d\n", totalTrials, numSignals)

	// パラメータの最適化を実行
	_, _, results := optimization.OptimizeParameters(&stockResponse, params)

	// 実行時間を測定
	elapsedTime := time.Since(startTime)

	// 結果を表示
	optimization.PrintResults(results, elapsedTime)

	// ここから各モデルの結果を表示するためのループを追加
	modelNames := []string{
		"LightGBM", "RandomForest", "XGBoost", "CatBoost", "AdaBoost",
		"DecisionTree", "GradientBoosting", "ExtraTrees", "Bagging", "Voting", "Stacking",
	}
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
			_, _, modelResults := optimization.OptimizeParameters(&stockResponse, params)

			// モデルの実行時間を測定
			modelElapsedTime := time.Since(modelStartTime)

			// モデルごとの結果を表示
			optimization.PrintResults(modelResults, modelElapsedTime, optimization.WithModelName(modelName), optimization.WithSignalCount(modelSignalCount))

			// 元のシグナルに戻す
			stockResponse.SymbolData[0].Signals = originalSignals
		} else {
			fmt.Printf("モデル: %s の予測データが見つかりませんでした。スキップします。\n", modelName)
		}
	}
}
