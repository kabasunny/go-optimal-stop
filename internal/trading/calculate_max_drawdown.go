package trading

// 後ほど、総資金をベースに計算を行う
// calculateMaxDrawdown は、トレード結果から最大ドローダウンを計算する
func calculateMaxDrawdown(tradeResults []tradeRecord) float64 {
	if len(tradeResults) == 0 {
		return 0
	}

	maxDrawdown := 0.0
	peak := 0.0             // 資産のピーク
	currentCapital := 100.0 // 仮に100スタート（任意の基準）

	for _, result := range tradeResults {
		currentCapital += result.ProfitLoss // 累積資産の計算
		if currentCapital > peak {
			peak = currentCapital
		}

		drawdown := (peak - currentCapital) / peak
		if drawdown > maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown * 100 // パーセンテージ表示
}
