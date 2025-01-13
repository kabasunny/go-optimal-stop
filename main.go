package main

import (
	"fmt"
	"os"
	"time"

	pb "go-optimal-stop/cmd/experiment_proto" // プロトコルバッファの定義をインポート
	"go-optimal-stop/internal/ml_stockdata"
	"go-optimal-stop/internal/optimization"

	"google.golang.org/protobuf/proto"
)

func main() {
	startTime := time.Now() // 実行時間の測定開始

	// 1570_2025-01-12.binを読み込み、stockResponseにプロトコルバッファバイナリからデータをマッピング
	filePath := "data/ml_stock_response/1570_2025-01-12.bin"
	data, err := os.ReadFile(filePath) // os.ReadFile に変更
	if err != nil {
		fmt.Printf("ファイルの読み込みエラー: %v\n", err)
		return
	}

	var protoResponse pb.MLStockResponse
	if err := proto.Unmarshal(data, &protoResponse); err != nil {
		fmt.Printf("プロトコルバッファのアンマーシャルエラー: %v\n", err)
		return
	}

	// プロトコルバッファから内部MLStockResponse型への変換
	stockResponse := convertProtoToInternal(&protoResponse)

	// Parameters構造体を作成し、関数を使ってパラメータを設定
	params := ml_stockdata.Parameters{}
	params.SetStopLoss(2.0, 5.0, 1.0)
	params.SetTrailingStop(5.0, 10.0, 1.0)
	params.SetTrailingStopUpdate(2.0, 5.0, 1.0)

	// protoResponse 内のシグナルの数を取得
	numSignals := len(stockResponse.SymbolData[0].Signals) // 最初のシンボルのシグナル数を使用

	// 総試行回数を算出
	totalTrials := len(params.StopLossPercentages) * len(params.TrailingStopTriggers) * len(params.TrailingStopUpdates) * len(stockResponse.SymbolData) * numSignals
	fmt.Printf("総試行回数: %d\n", totalTrials)

	// パラメータの最適化を実行
	bestResult, worstResult, _ := optimization.OptimizeParameters(&stockResponse, params)

	// 実行時間を測定
	elapsedTime := time.Since(startTime)

	// 結果を表示
	fmt.Printf("最良の結果: {StopLossPercentage: %.2f%%, TrailingStopTrigger: %.2f%%, TrailingStopUpdate: %.2f%%, ProfitLoss: %.2f%%, PurchaseDate: %s, ExitDate: %s}\n",
		bestResult.StopLossPercentage, bestResult.TrailingStopTrigger, bestResult.TrailingStopUpdate, bestResult.ProfitLoss, bestResult.PurchaseDate, bestResult.ExitDate)
	fmt.Printf("最悪の結果: {StopLossPercentage: %.2f%%, TrailingStopTrigger: %.2f%%, TrailingStopUpdate: %.2f%%, ProfitLoss: %.2f%%, PurchaseDate: %s, ExitDate: %s}\n",
		worstResult.StopLossPercentage, worstResult.TrailingStopTrigger, worstResult.TrailingStopUpdate, worstResult.ProfitLoss, worstResult.PurchaseDate, worstResult.ExitDate)
	fmt.Printf("実行時間: %v\n", elapsedTime)
}

// プロトコルバッファから内部MLStockResponse型への変換関数
func convertProtoToInternal(protoResponse *pb.MLStockResponse) ml_stockdata.MLStockResponse {
	var stockResponse ml_stockdata.MLStockResponse
	for _, protoSymbolData := range protoResponse.SymbolData {
		var symbolData ml_stockdata.MLSymbolData
		symbolData.Symbol = protoSymbolData.Symbol
		symbolData.Signals = protoSymbolData.Signals
		for _, protoDailyData := range protoSymbolData.DailyData {
			dailyData := ml_stockdata.MLDailyData{
				Date:   protoDailyData.GetDate(),           // GetXXXメソッドを使用
				Open:   float64(protoDailyData.GetOpen()),  // 型変換を追加
				High:   float64(protoDailyData.GetHigh()),  // 型変換を追加
				Low:    float64(protoDailyData.GetLow()),   // 型変換を追加
				Close:  float64(protoDailyData.GetClose()), // 型変換を追加
				Volume: protoDailyData.GetVolume(),
			}
			symbolData.DailyData = append(symbolData.DailyData, dailyData)
		}
		stockResponse.SymbolData = append(stockResponse.SymbolData, symbolData)
	}
	return stockResponse
}
