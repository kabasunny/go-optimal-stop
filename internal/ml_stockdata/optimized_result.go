// internal/stockdata/result.go

package ml_stockdata

// OptimizedResult 構造体の定義
type OptimizedResult struct {
	StopLossPercentage  float64
	TrailingStopTrigger float64
	TrailingStopUpdate  float64
	ProfitLoss          float64 // 損益率
	WinRate             float64 // 勝率
	MaxProfit           float64 // 連続して加算された利益の幅
	MaxLoss             float64 // 連続して加算された損失の幅
	PurchaseDate        string
	ExitDate            string
}
