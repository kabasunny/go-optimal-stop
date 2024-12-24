// internal/data/data.go

package data

import "time"

// Data 構造体の定義
type Data struct {
	Date  time.Time
	Open  float64
	Low   float64
	Close float64
}

// Result 構造体の定義
type Result struct {
	StopLossPercentage  float64
	TrailingStopTrigger float64
	TrailingStopUpdate  float64
	ProfitLoss          float64
	PurchaseDate        time.Time
	ExitDate            time.Time
}
