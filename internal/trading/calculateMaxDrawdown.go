package trading

// calculateMaxDrawdown 関数は、トレード結果から最大ドローダウンを計算
func calculateMaxDrawdown(tradeResults []tradeResult) float64 {
	maxDrawdown := 0.0
	peak := tradeResults[0].ProfitLoss

	for _, result := range tradeResults {
		if result.ProfitLoss > peak {
			peak = result.ProfitLoss
		}
		drawdown := (peak - result.ProfitLoss) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown * 100 // パーセンテージ表示
}
