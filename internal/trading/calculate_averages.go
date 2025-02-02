package trading

func calculateAverages(tradeResults []tradeRecord) (float64, float64) {
	var totalProfit, totalLoss float64
	winCount, lossCount := 0, 0

	for _, result := range tradeResults {
		if result.ProfitLoss > 0 {
			totalProfit += result.ProfitLoss
			winCount++
		} else {
			totalLoss += result.ProfitLoss
			lossCount++
		}
	}

	averageProfit := 0.0
	averageLoss := 0.0
	if winCount > 0 {
		averageProfit = totalProfit / float64(winCount)
	}
	if lossCount > 0 {
		averageLoss = totalLoss / float64(lossCount)
	}

	return averageProfit, averageLoss
}
