package trading

// calculateDrawdownAndDrawup は、トレード結果から最大ドローダウンと最大上昇率を計算する
func calculateDrawdownAndDrawup(tradeResults []tradeRecord) (float64, float64) {
	if len(tradeResults) == 0 {
		return 0, 0
	}

	maxDrawdown := 0.0
	maxDrawup := 0.0
	peak := tradeResults[0].PortfolioValue   // 初期資産をピークとして設定
	trough := tradeResults[0].PortfolioValue // 初期資産をトラフとして設定

	for _, result := range tradeResults {
		// ドローダウンの計算
		if result.PortfolioValue > peak {
			peak = result.PortfolioValue
		}
		drawdown := float64(peak-result.PortfolioValue) / float64(peak)
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}

		// 最大上昇率の計算
		if result.PortfolioValue < trough {
			trough = result.PortfolioValue
		}
		drawup := float64(result.PortfolioValue-trough) / float64(trough)
		if drawup > maxDrawup {
			maxDrawup = drawup
		}
	}

	return maxDrawdown * 100, maxDrawup * 100 // パーセンテージ表示
}
