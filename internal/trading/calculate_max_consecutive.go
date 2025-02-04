package trading

// 連続利益と損失を計算する関数
func calculateMaxConsecutive(trades []tradeRecord) (float64, float64) {
	maxConsecutiveProfit := 0.0
	maxConsecutiveLoss := 0.0
	currentConsecutiveProfit := 0.0
	currentConsecutiveLoss := 0.0

	for _, trade := range trades {
		// 金額ベースの利益・損失を計算
		profitLossAmount := (trade.ProfitLoss / 100) * trade.EntryCost

		if profitLossAmount > 0 {
			// 利益がある場合、連続利益に加算
			currentConsecutiveProfit += profitLossAmount
			if currentConsecutiveProfit > maxConsecutiveProfit {
				maxConsecutiveProfit = currentConsecutiveProfit
			}
			// 連続損失をリセット
			currentConsecutiveLoss = 0
		} else {
			// 損失がある場合、連続損失に加算
			currentConsecutiveLoss += profitLossAmount
			if currentConsecutiveLoss < maxConsecutiveLoss {
				maxConsecutiveLoss = currentConsecutiveLoss
			}
			// 連続利益をリセット
			currentConsecutiveProfit = 0
		}
	}

	// 使わないけど、金額を返す様に変更
	return maxConsecutiveProfit, maxConsecutiveLoss
}
