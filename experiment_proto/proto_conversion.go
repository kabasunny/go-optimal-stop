package experiment_proto

import (
	"go-optimal-stop/internal/ml_stockdata"
)

func ConvertProtoToInternal(protoResponse *MLStockResponse) ml_stockdata.InMLStockResponse {
	var stockResponse ml_stockdata.InMLStockResponse
	for _, protoSymbolData := range protoResponse.SymbolData {
		var symbolData ml_stockdata.InMLSymbolData
		symbolData.Symbol = protoSymbolData.Symbol
		symbolData.Signals = protoSymbolData.Signals
		for _, protoDailyData := range protoSymbolData.DailyData {
			dailyData := ml_stockdata.InMLDailyData{
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
