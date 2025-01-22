package trading

import (
	"time"
)

// 並列処理された取引結果を収集し、ソート、連続ストリークの計算などの処理を安全に行うために使用
type tradeResult struct {
	Symbol     string
	Date       time.Time
	ProfitLoss float64
}
