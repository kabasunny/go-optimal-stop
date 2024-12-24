// internal/stockdata/symbol_data.go

package stockdata

// SymbolData 構造体の定義
type SymbolData struct {
	Symbol    string   // 銘柄コード
	DailyData []Data   // 株価データ
	Signals   []string // 日付のリスト
}
