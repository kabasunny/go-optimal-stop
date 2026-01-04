package experiment_proto

import (
	"go-optimal-stop/internal/ml_stockdata"
)

// ConvertProtoToInternal は、MLStockResponse を InMLStockResponse に変換する関数です
func ConvertProtoToInternal(protoResponse *MLStockResponse) ml_stockdata.InMLStockResponse {
	var stockResponse ml_stockdata.InMLStockResponse
	for _, protoSymbolData := range protoResponse.SymbolData {
		var symbolData ml_stockdata.InMLSymbolData
		symbolData.Symbol = protoSymbolData.Symbol

		// 文字列のシグナルをSignal構造体に変換（既存データとの互換性のため、すべて買いシグナルとして扱う）
		for _, signalDate := range protoSymbolData.Signals {
			signal := ml_stockdata.Signal{
				Date:       signalDate,
				Type:       ml_stockdata.SignalTypeBuy, // デフォルトは買いシグナル
				Confidence: 1.0,                        // デフォルトの信頼度
			}
			symbolData.Signals = append(symbolData.Signals, signal)
		}

		symbolData.Priority = protoSymbolData.Priority // int64として直接割り当て
		for _, protoDailyData := range protoSymbolData.DailyData {
			dailyData := ml_stockdata.InMLDailyData{
				Date:   protoDailyData.GetDate(),
				Open:   float64(protoDailyData.GetOpen()),
				High:   float64(protoDailyData.GetHigh()),
				Low:    float64(protoDailyData.GetLow()),
				Close:  float64(protoDailyData.GetClose()),
				Volume: protoDailyData.GetVolume(),
			}
			symbolData.DailyData = append(symbolData.DailyData, dailyData)
		}
		stockResponse.SymbolData = append(stockResponse.SymbolData, symbolData)
	}
	return stockResponse
}
