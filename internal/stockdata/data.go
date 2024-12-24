// internal/stockdata/data.go

package stockdata

// Data 構造体の定義
type Data struct {
	Date  string // 日付を文字列として保持
	Open  float64
	Low   float64
	Close float64
}
