// internal/stockdata/result.go

package ml_stockdata

// OptimizedResult 構造体の定義
type OptimizedResult struct {
	StopLossPercentage  float64
	TrailingStopTrigger float64
	TrailingStopUpdate  float64
	ProfitLoss          float64
	WinRate             float64
	PurchaseDate        string
	ExitDate            string
}
