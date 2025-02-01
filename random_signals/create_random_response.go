package random_signals

import (
	"fmt"
	"os"

	"go-optimal-stop/experiment_proto"
	"go-optimal-stop/internal/ml_stockdata"

	"google.golang.org/protobuf/proto"
)

// CSVファイルからデータを読み込み、StockResponse構造体を作成
func createStockResponse(filePath *string, seed ...int64) (ml_stockdata.InMLStockResponse, int, []string, error) {
	// ファイルを読み込み、stockResponseにプロトコルバッファバイナリからデータをマッピング
	data, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("ファイルの読み込みエラー: %v\n", err)
		return ml_stockdata.InMLStockResponse{}, 0, nil, err
	}

	var protoResponse experiment_proto.MLStockResponse
	if err := proto.Unmarshal(data, &protoResponse); err != nil {
		fmt.Printf("プロトコルバッファのアンマーシャルエラー: %v\n", err)
		return ml_stockdata.InMLStockResponse{}, 0, nil, err
	}

	// プロトコルバッファから内部MLStockResponse型への変換
	stockResponse := experiment_proto.ConvertProtoToInternal(&protoResponse)

	// 抽出したシンボル
	var symbols []string

	// 各銘柄のシグナル数を合計する
	totalSignals := 0
	for i := range stockResponse.SymbolData {
		numSignals := len(stockResponse.SymbolData[i].Signals)
		symbols = append(symbols, stockResponse.SymbolData[i].Symbol)

		// 検出したシグナル数でランダムにシグナルを生成
		var signals []string
		if len(seed) > 0 {
			signals = generateRandomSignals(stockResponse.SymbolData[i].DailyData, numSignals, seed[0])
		} else {
			signals = generateRandomSignals(stockResponse.SymbolData[i].DailyData, numSignals)
		}
		stockResponse.SymbolData[i].Signals = signals

		totalSignals += numSignals
	}

	return stockResponse, totalSignals, symbols, nil
}
