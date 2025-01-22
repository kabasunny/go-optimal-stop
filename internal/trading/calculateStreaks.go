package trading

// calculateStreaks 関数は、トレード結果のスライスを受け取り、最大ポジティブストリークと最大ネガティブストリークを計算して返す
func calculateStreaks(results []tradeResult) (float64, float64) {
	maxPositiveStreak := 0.0     // 最大のポジティブストリークを追跡
	maxNegativeStreak := 0.0     // 最大のネガティブストリークを追跡
	currentPositiveStreak := 0.0 // 現在のポジティブストリークを追跡
	currentNegativeStreak := 0.0 // 現在のネガティブストリークを追跡

	// 各トレード結果をループ処理
	for _, result := range results {
		// 利益が正の値の場合
		if result.ProfitLoss > 0 {
			currentPositiveStreak += result.ProfitLoss // 現在のポジティブストリークに利益を加算
			if currentPositiveStreak > maxPositiveStreak {
				maxPositiveStreak = currentPositiveStreak // 必要に応じて最大ポジティブストリークを更新
			}
			currentNegativeStreak = 0 // ネガティブストリークをリセット
		} else {
			// 利益が負の値の場合
			currentNegativeStreak += result.ProfitLoss // 現在のネガティブストリークに損失を加算
			if currentNegativeStreak < maxNegativeStreak {
				maxNegativeStreak = currentNegativeStreak // 必要に応じて最大ネガティブストリークを更新
			}
			currentPositiveStreak = 0 // ポジティブストリークをリセット
		}
	}

	// 最大ポジティブストリークと最大ネガティブストリークを返す
	return maxPositiveStreak, maxNegativeStreak
}
