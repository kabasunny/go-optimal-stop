package random_signals

import (
	"fmt"

	"go-optimal-stop/internal/ml_stockdata"
)

// CSVファイルからデータを読み込み、StockResponse構造体を作成
func createStockResponse(csvDir string, symbols []string, numSignals int, startDate string, seed ...int64) (ml_stockdata.InMLStockResponse, error) {
	var symbolDataList []ml_stockdata.InMLSymbolData
	var allDailyData []ml_stockdata.InMLDailyData

	for _, symbol := range symbols {
		filePath := fmt.Sprintf("%s/%s.csv", csvDir, symbol)

		// CSVファイルを読み込み
		data, err := loadCSV(filePath, startDate)
		if err != nil {
			return ml_stockdata.InMLStockResponse{}, fmt.Errorf("CSVファイルの読み込みエラー: %v", err)
		}

		symbolData := ml_stockdata.InMLSymbolData{
			Symbol:    symbol,
			DailyData: data,
		}
		symbolDataList = append(symbolDataList, symbolData)
		allDailyData = append(allDailyData, data...)
	}

	// ランダムにシグナルを生成
	var signals []string
	if len(seed) > 0 {
		signals = generateRandomSignals(allDailyData, numSignals, seed[0])
	} else {
		signals = generateRandomSignals(allDailyData, numSignals)
	}

	// シグナルをシンボルごとに分割してマッピング
	signalsPerSymbol := numSignals / len(symbols)
	for i := range symbolDataList {
		startIndex := i * signalsPerSymbol
		endIndex := startIndex + signalsPerSymbol
		if endIndex > len(signals) {
			endIndex = len(signals)
		}
		symbolDataList[i].Signals = signals[startIndex:endIndex]
	}

	// StockResponse構造体を作成
	return ml_stockdata.InMLStockResponse{
		SymbolData: symbolDataList,
	}, nil
}
