package trading

import (
	"time"
)

// tradeRecord 構造体は取引結果を収集し、ソート、連続ストリークの計算などの処理を安全に行うために使用
type tradeRecord struct {
	Symbol       string
	EntryDate    time.Time
	ExitDate     time.Time
	ProfitLoss   float64 // ％
	EntryCost    float64 // 一株の価格 × 株数 + 手数料
	PositionSize float64 // 株数
	EntryPrice   float64 // 一株の価格
	ExitPrice    float64 // 一株の価格
}
