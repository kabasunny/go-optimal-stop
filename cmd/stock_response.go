// cmd/stock_response.go

package main

import (
	"fmt"

	"go-optimal-stop/internal/stockdata"
)

// CSVファイルからデータを読み込み、StockResponse構造体を作成
func createStockResponse(csvDir string, symbols []string, numSignals int, seed ...int64) (stockdata.StockResponse, error) {
	var symbolDataList []stockdata.SymbolData

	for _, symbol := range symbols {
		filePath := fmt.Sprintf("%s/%s_stock_data.csv", csvDir, symbol)

		// CSVファイルを読み込み
		data, err := loadCSV(filePath)
		if err != nil {
			return stockdata.StockResponse{}, fmt.Errorf("CSVファイルの読み込みエラー: %v", err)
		}

		// ランダムにシグナルを生成
		var signals []string
		if len(seed) > 0 {
			signals = generateRandomSignals(data, numSignals, seed[0])
		} else {
			signals = generateRandomSignals(data, numSignals)
		}

		// SymbolData構造体を作成
		symbolData := stockdata.SymbolData{
			Symbol:    symbol,
			DailyData: data,
			Signals:   signals,
		}
		symbolDataList = append(symbolDataList, symbolData)
	}

	// StockResponse構造体を作成
	return stockdata.StockResponse{
		SymbolData: symbolDataList,
	}, nil
}
