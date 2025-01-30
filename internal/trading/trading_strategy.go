package trading

import (
	"go-optimal-stop/internal/ml_stockdata"
	"sort"
	"time"
)

// TradingStrategy 関数は、与えられた株価データとトレーディングパラメータに基づいて総利益、勝率、最大ポジティブストリーク、最大ネガティブストリークを返す
func TradingStrategy(response *ml_stockdata.InMLStockResponse, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (float64, float64, float64, float64, error) {
	totalProfitLoss := 0.0         // 全体の利益を追跡
	winCount := 0                  // 勝ちトレードのカウント
	totalCount := 0                // 全トレードのカウント
	var tradeResults []tradeResult // トレード結果を保持するスライス

	// continue した日付を記録するリスト
	// var skippedDates []time.Time
	var totalSignals int

	// 各シンボルデータをループ処理
	for _, symbolData := range response.SymbolData {
		previousEndDate := time.Time{} // 前回の終了日を記録する変数

		// 各シグナルをループ処理
		for _, signal := range symbolData.Signals {
			totalSignals++ // シグナルの総数をカウント

			startDate, err := parseDate(signal) // シグナルの日付を解析
			if err != nil {
				// skippedDates = append(skippedDates, startDate)
				continue
			}

			// 前回の終了日と開始日が重なる場合、次の開始日に移る
			if !previousEndDate.IsZero() && startDate.Before(previousEndDate) {
				// skippedDates = append(skippedDates, startDate)
				continue
			}

			// トレード戦略を実行し、利益を計算
			purchaseDate, endDate, profitLoss, _, _, err := singleTradingStrategy(&symbolData.DailyData, startDate, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
			if err != nil {
				// skippedDates = append(skippedDates, startDate)
				continue
			}

			totalProfitLoss += profitLoss // 総利益に加算
			totalCount++                  // トレード数をインクリメント
			if profitLoss > 0 {
				winCount++ // 勝ちトレードの場合、勝ち数をインクリメント
			}

			// トレード結果をスライスに追加
			tradeResults = append(tradeResults, tradeResult{
				Symbol:     symbolData.Symbol,
				Date:       purchaseDate,
				ProfitLoss: profitLoss,
			})

			// 前回の終了日を更新
			previousEndDate = endDate
		}
	}

	// トレード結果をシンボルと日付でソート
	sort.Slice(tradeResults, func(i, j int) bool {
		if tradeResults[i].Symbol == tradeResults[j].Symbol {
			return tradeResults[i].Date.Before(tradeResults[j].Date)
		}
		return tradeResults[i].Symbol < tradeResults[j].Symbol
	})

	// 最大ポジティブストリークと最大ネガティブストリークを計算
	maxPositiveStreak, maxNegativeStreak := calculateStreaks(tradeResults)

	// 勝率を計算
	winRate := float64(winCount) / float64(totalCount) * 100

	// スキップしたシグナルの個数を表示
	// skippedCount := len(skippedDates)
	// fmt.Printf("Total Signals: %d, Skipped Signals: %d, Processed Signals: %d\n", totalSignals, skippedCount, totalSignals-skippedCount)

	return totalProfitLoss, winRate, maxPositiveStreak, maxNegativeStreak, nil
}
