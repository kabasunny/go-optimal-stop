package trading

import (
	"time"
)

// tradeRecord 構造体は取引結果を収集し、ソート、連続ストリークの計算などの処理を安全に行うために使用
type tradeRecord struct {
	Symbol       string
	EntryDate    time.Time
	ExitDate     time.Time
	ProfitLoss   float64
	EntryCost    float64
	ExitPrice    float64
	PositionSize float64
	EntryPrice   float64
}
