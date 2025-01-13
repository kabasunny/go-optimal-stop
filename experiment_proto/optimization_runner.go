package experiment_proto

import (
	"fmt"
	"os"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/optimization"

	"google.golang.org/protobuf/proto"
)

func RunOptimization(filePath string) {
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

	// Parameters構造体を作成し、関数を使ってパラメータを設定
	params := ml_stockdata.Parameters{}
	params.SetStopLoss(2.0, 5.0, 1.0)
	params.SetTrailingStop(5.0, 10.0, 1.0)
	params.SetTrailingStopUpdate(2.0, 5.0, 1.0)

	// protoResponse 内のシグナルの数を取得
	numSignals := len(stockResponse.SymbolData[0].Signals) // 最初のシンボルのシグナル数を使用

	// 総試行回数を算出
	totalTrials := len(params.StopLossPercentages) * len(params.TrailingStopTriggers) * len(params.TrailingStopUpdates) * len(stockResponse.SymbolData) * numSignals
	fmt.Printf("総試行回数: %d, シグナル数: %d\n", totalTrials, numSignals)

	// パラメータの最適化を実行
	_, _, results := optimization.OptimizeParameters(&stockResponse, params)

	// 実行時間を測定
	elapsedTime := time.Since(startTime)

	// 結果を表示
	optimization.PrintOverallResults(results, elapsedTime)

	// 各モデルの結果を表示
	modelNames := []string{"LightGBM", "RandomForest", "XGBoost", "CatBoost", "AdaBoost", "SVM", "KNeighbors", "LogisticRegression"}
	for _, modelName := range modelNames {
		fmt.Printf("モデル: %s\n", modelName)
		modelSignals := protoResponse.SymbolData[0].ModelPredictions[modelName].PredictionDates
		stockResponse.SymbolData[0].Signals = modelSignals // モデルごとのシグナルを設定

		// シグナル数を取得
		modelSignalCount := len(modelSignals)

		// モデルの最適化開始時間を記録
		modelStartTime := time.Now()

		// パラメータの最適化を再実行
		_, _, modelResults := optimization.OptimizeParameters(&stockResponse, params)

		// モデルの実行時間を測定
		modelElapsedTime := time.Since(modelStartTime)

		// モデルごとの結果を表示
		optimization.PrintModelResults(modelName, modelSignalCount, modelResults, modelElapsedTime)
	}
}
