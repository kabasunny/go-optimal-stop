// internal/stockdata/result.go

package ml_stockdata

// OptimizedionResult 構造体の定義
type OptimizedionResult struct {
	StopLossPercentage  float64
	TrailingStopTrigger float64
	TrailingStopUpdate  float64
	ProfitLoss          float64
	PurchaseDate        string
	ExitDate            string
}
