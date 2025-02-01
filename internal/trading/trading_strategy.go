package trading

import (
	"go-optimal-stop/internal/ml_stockdata"
	"math"
	"sort"
	"time"
)

// TradingStrategy 関数は、与えられた株価データとトレーディングパラメータに基づいて総利益、勝率、最大ポジティブストリーク、最大ネガティブストリーク、その他の指標を返す
func TradingStrategy(response *ml_stockdata.InMLStockResponse, totalFunds *int, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (float64, float64, float64, float64, int, int, float64, float64, float64, float64, float64, float64, error) {

	// ここで総資金の変数を用意する、（前提：どの銘柄のシグナルを優先するか）
	// currentTotalFunds := *totalFunds

	totalProfitLoss := 0.0         // 全体の利益を追跡
	winCount := 0                  // 勝ちトレードのカウント
	totalCount := 0                // 全トレードのカウント
	var tradeResults []tradeResult // トレード結果を保持するスライス

	// 各シンボルデータをループ処理
	for _, symbolData := range response.SymbolData {
		previousEndDate := time.Time{} // 前回の終了日を記録する変数

		// 各シグナルをループ処理
		for _, signal := range symbolData.Signals {
			startDate, err := parseDate(signal) // シグナルの日付を解析
			if err != nil {
				continue
			}

			// 前回の終了日と開始日が重なる場合、次の開始日に移る
			if !previousEndDate.IsZero() && startDate.Before(previousEndDate) {
				continue
			}

			// トレード戦略を実行し、利益を計算
			purchaseDate, endDate, profitLoss, _, _, err := singleTradingStrategy(&symbolData.DailyData, startDate, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
			if err != nil {
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

	// 平均利益率と平均損失率を計算
	var totalProfit, totalLoss float64
	if winCount > 0 {
		for _, result := range tradeResults {
			if result.ProfitLoss > 0 {
				totalProfit += result.ProfitLoss
			} else {
				totalLoss += result.ProfitLoss
			}
		}
	}
	averageProfit := 0.0
	averageLoss := 0.0
	if winCount > 0 {
		averageProfit = totalProfit / float64(winCount)
	}
	if totalCount-winCount > 0 {
		averageLoss = totalLoss / float64(totalCount-winCount)
	}

	// 最大ドローダウンを計算
	// 後ほど、総資金をベースに計算を行う
	maxDrawdown := calculateMaxDrawdown(tradeResults)

	// 超過リターンを計算
	excessReturns := []float64{}
	for _, result := range tradeResults {
		excessReturns = append(excessReturns, result.ProfitLoss)
	}

	// シャープレシオを計算（リスクフリーレートを0と仮定）
	sharpeRatio := 0.0
	if stdDev := standardDeviation(excessReturns); stdDev > 0 {
		sharpeRatio = calculateSharpeRatio(tradeResults, 0)
	}

	// リスクリワード比を計算
	riskRewardRatio := 0.0
	if averageLoss != 0 {
		riskRewardRatio = averageProfit / math.Abs(averageLoss)
	}

	// 期待値を計算（パーセンテージ表示）
	expectedValue := 0.0
	if totalCount > 0 {
		expectedValue = ((winRate * averageProfit) - ((100 - winRate) * averageLoss)) / 100
	}

	return totalProfitLoss, winRate, maxPositiveStreak, maxNegativeStreak, winCount, totalCount - winCount, averageProfit, averageLoss, maxDrawdown, sharpeRatio, riskRewardRatio, expectedValue, nil
}
