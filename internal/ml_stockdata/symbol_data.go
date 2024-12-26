// internal/stockdata/symbol_data.go

package ml_stockdata

// MLSymbolData 構造体の定義
type MLSymbolData struct {
	Symbol    string   // 銘柄コード
	DailyData []Data   // 株価データ
	Signals   []string // 日付のリスト
}
