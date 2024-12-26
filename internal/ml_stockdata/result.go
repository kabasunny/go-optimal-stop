// internal/stockdata/result.go

package ml_stockdata

// Result 構造体の定義
type Result struct {
	StopLossPercentage  float64
	TrailingStopTrigger float64
	TrailingStopUpdate  float64
	ProfitLoss          float64
	PurchaseDate        string
	ExitDate            string
}
