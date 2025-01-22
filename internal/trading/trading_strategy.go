package trading

import (
	"sort"
	"sync"

	"go-optimal-stop/internal/ml_stockdata"
)

// TradingStrategy 関数は、与えられた株価データとトレーディングパラメータに基づいて総利益、勝率、最大ポジティブストリーク、最大ネガティブストリークを返す
func TradingStrategy(response *ml_stockdata.InMLStockResponse, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (float64, float64, float64, float64, error) {
	totalProfitLoss := 0.0         // 全体の利益を追跡
	winCount := 0                  // 勝ちトレードのカウント
	totalCount := 0                // 全トレードのカウント
	var tradeResults []tradeResult // トレード結果を保持するスライス
	var mu sync.Mutex              // 排他制御用のミューテックス
	var wg sync.WaitGroup          // 同期用のWaitGroup

	// 各シンボルデータをループ処理
	for _, symbolData := range response.SymbolData {
		// 各シグナルをループ処理
		for _, signal := range symbolData.Signals {
			wg.Add(1) // WaitGroupのカウントをインクリメント
			go func(symbolData ml_stockdata.InMLSymbolData, signal string) {
				defer wg.Done()                     // 処理終了時にWaitGroupのカウントをデクリメント
				startDate, err := parseDate(signal) // シグナルの日付を解析
				if err != nil {
					return
				}

				// トレード戦略を実行し、利益を計算
				_, _, profitLoss, _, _, err := singleTradingStrategy(&symbolData.DailyData, startDate, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
				if err != nil {
					return
				}
				mu.Lock()                     // 排他制御開始
				totalProfitLoss += profitLoss // 総利益に加算
				totalCount++                  // トレード数をインクリメント
				if profitLoss > 0 {
					winCount++ // 勝ちトレードの場合、勝ち数をインクリメント
				}
				// トレード結果をスライスに追加
				tradeResults = append(tradeResults, tradeResult{
					Symbol:     symbolData.Symbol,
					Date:       startDate,
					ProfitLoss: profitLoss,
				})
				mu.Unlock() // 排他制御終了
			}(symbolData, signal)
		}
	}

	wg.Wait() // すべてのGoルーチンの終了を待機

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
	return totalProfitLoss, winRate, maxPositiveStreak, maxNegativeStreak, nil
}
