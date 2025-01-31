package trading

// calculateSharpeRatio 関数は、トレード結果からシャープレシオを計算
// シャープレシオはリスク調整後のリターンを評価する指標
func calculateSharpeRatio(tradeResults []tradeResult, riskFreeRate float64) float64 {
	totalReturn := 0.0
	returns := []float64{}

	// 各トレード結果の利益を収集し、総利益を計算
	for _, result := range tradeResults {
		returns = append(returns, result.ProfitLoss)
		totalReturn += result.ProfitLoss
	}

	// 平均リターンを計算
	meanReturn := totalReturn / float64(len(returns))
	excessReturns := []float64{}

	// 各リターンからリスクフリーレートを引いた超過リターンを計算
	for _, r := range returns {
		excessReturns = append(excessReturns, r-riskFreeRate)
	}

	// 超過リターンの標準偏差を計算
	stdDev := standardDeviation(excessReturns)

	// シャープレシオを計算
	return (meanReturn - riskFreeRate) / stdDev
}
