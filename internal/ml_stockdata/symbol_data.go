package ml_stockdata

// SignalType シグナルの種類を表す列挙型
type SignalType int

const (
	SignalTypeBuy SignalType = iota
	SignalTypeSell
)

// Signal 買い・売りシグナルの情報
type Signal struct {
	Date       string
	Type       SignalType
	Confidence float64
}

// InMLSymbolData 構造体の定義
type InMLSymbolData struct {
	Symbol    string          // 銘柄コード
	DailyData []InMLDailyData // 株価データ
	Signals   []Signal        // 買い・売りシグナルのリスト（変更：[]string から []Signal へ）
	Priority  int64           // 優先順位
}
