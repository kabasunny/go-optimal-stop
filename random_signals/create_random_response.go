package random_signals

import (
	"fmt"

	"go-optimal-stop/internal/ml_stockdata"
)

// CSVファイルからデータを読み込み、StockResponse構造体を作成
func createStockResponse(csvDir string, symbols []string, numSignals int, seed ...int64) (ml_stockdata.InMLStockResponse, error) {
	var symbolDataList []ml_stockdata.InMLSymbolData

	for _, symbol := range symbols {
		filePath := fmt.Sprintf("%s/%s.csv", csvDir, symbol)

		// CSVファイルを読み込み
		data, err := loadCSV(filePath)

		if err != nil {
			return ml_stockdata.InMLStockResponse{}, fmt.Errorf("CSVファイルの読み込みエラー: %v", err)
		}

		// ランダムにシグナルを生成
		var signals []string
		if len(seed) > 0 {
			signals = generateRandomSignals(data, numSignals, seed[0])
		} else {
			signals = generateRandomSignals(data, numSignals)
		}

		// SymbolData構造体を作成
		symbolData := ml_stockdata.InMLSymbolData{
			Symbol:    symbol,
			DailyData: data,
			Signals:   signals,
		}
		symbolDataList = append(symbolDataList, symbolData)
	}

	// StockResponse構造体を作成
	return ml_stockdata.InMLStockResponse{
		SymbolData: symbolDataList,
	}, nil
}
