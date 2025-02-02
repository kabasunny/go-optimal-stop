package ml_stockdata

// InMLSymbolData 構造体の定義
type InMLSymbolData struct {
	Symbol    string          // 銘柄コード
	DailyData []InMLDailyData // 株価データ
	Signals   []string        // 日付のリスト
	Priority  int64           // 優先順位
}
