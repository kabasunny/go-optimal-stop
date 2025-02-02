package trading

// 連続利益と損失を計算する関数
func calculateMaxConsecutive(trades []tradeRecord) (float64, float64) {
	maxConsecutiveProfit := 0.0
	maxConsecutiveLos := 0.0
	currentConsecutiveProfit := 0.0
	currentConsecutiveLos := 0.0

	for _, trade := range trades {
		if trade.ProfitLoss > 0 {
			currentConsecutiveProfit += trade.ProfitLoss
			if currentConsecutiveProfit > maxConsecutiveProfit {
				maxConsecutiveProfit = currentConsecutiveProfit
			}
			currentConsecutiveLos = 0
		} else {
			currentConsecutiveLos += trade.ProfitLoss
			if currentConsecutiveLos < maxConsecutiveLos {
				maxConsecutiveLos = currentConsecutiveLos
			}
			currentConsecutiveProfit = 0
		}
	}

	return maxConsecutiveProfit, maxConsecutiveLos
}
