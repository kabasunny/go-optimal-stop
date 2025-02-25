// internal/stockdata/data.go

package ml_stockdata

// InMLDailyData 構造体の定義
type InMLDailyData struct {
	Date   string // 日付を文字列として保持
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}
