package main

import (
	"fmt"
	"go-optimal-stop/internal/ml_stockdata"
)

// CSVファイルからデータを読み込み、StockResponse構造体を作成
func CreateStockResponse(csvDir string, symbols []string, numSignals int, seed ...int64) (ml_stockdata.MLStockResponse, error) {
	var symbolDataList []ml_stockdata.MLSymbolData

	for _, symbol := range symbols {
		filePath := fmt.Sprintf("%s/%s_stock_data.csv", csvDir, symbol)

		// CSVファイルを読み込み
		data, err := LoadCSV(filePath)
		if err != nil {
			return ml_stockdata.MLStockResponse{}, fmt.Errorf("CSVファイルの読み込みエラー: %v", err)
		}

		// ランダムにシグナルを生成
		var signals []string
		if len(seed) > 0 {
			signals = generateRandomSignals(data, numSignals, seed[0])
		} else {
			signals = generateRandomSignals(data, numSignals)
		}

		// SymbolData構造体を作成
		symbolData := ml_stockdata.MLSymbolData{
			Symbol:    symbol,
			DailyData: data,
			Signals:   signals,
		}
		symbolDataList = append(symbolDataList, symbolData)
	}

	// StockResponse構造体を作成
	return ml_stockdata.MLStockResponse{
		SymbolData: symbolDataList,
	}, nil
}
