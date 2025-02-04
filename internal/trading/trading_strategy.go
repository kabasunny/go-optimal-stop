package trading

import (
	"fmt"
	"go-optimal-stop/internal/ml_stockdata"
	"math"
	"sort"
	"time"
)

// / TradingStrategy 関数は、与えられた株価データとトレーディングパラメータに基づいて総利益、勝率、その他の指標を返す
func TradingStrategy(response *ml_stockdata.InMLStockResponse, totalFunds *int, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (float64, float64, float64, float64, int, int, float64, float64, float64, float64, float64, float64, error) {

	// エントリー可能金額までのエントリー順序を決定する
	// まずは空のスライスを作成
	signals := []struct {
		Symbol     string
		SignalDate time.Time
		Priority   int64
	}{}

	// 各銘柄のシグナルを取得し、日付順にソート
	for _, symbolData := range response.SymbolData {
		if len(symbolData.Signals) < 1 {
			// fmt.Println("シグナルがないためスキップします")
			continue
		}
		for _, signal := range symbolData.Signals {
			date, err := parseDate(signal)
			if err != nil {
				fmt.Println("signal skip")
				continue // 日付の解析に失敗した場合はスキップ
			}
			// シグナル情報を追加
			signals = append(signals, struct {
				Symbol     string
				SignalDate time.Time
				Priority   int64
			}{symbolData.Symbol, date, symbolData.Priority})
		}
	}

	// シグナルを日付順、優先順にソート
	sort.Slice(signals, func(i, j int) bool { // インデックス i と j の要素を比較し、i が j よりも前に来るべき → true
		if signals[i].SignalDate.Equal(signals[j].SignalDate) {
			return signals[i].Priority < signals[j].Priority // シグナル日付が同じ場合、Priorityが小さい方が優先
		}
		return signals[i].SignalDate.Before(signals[j].SignalDate)
	})
	// fmt.Print(signals)

	activeTrades := make(map[string]tradeRecord) // 各シンボルのホールド状態
	originalTotalFunds := *totalFunds            // 総資金の初期化（コピーを作成）
	portfolioValue := originalTotalFunds         // ポートフォリオ額
	availableFunds := portfolioValue             // 使用可能な資金の初期化
	totalProfitLoss := 0.0                       // 全体の利益を追跡
	winCount, totalCount := 0, 0                 // 勝ちトレード数と総トレード数
	var tradeResults []tradeRecord               // トレード結果を保持するスライス

	// シンボルごとのエグジット情報を保持するマップ
	exitMap := make(map[time.Time][]tradeRecord)

	// ---- シグナルの処理 ----
	for _, signal := range signals {
		// fmt.Println("シグナル処理中:", signal) // デバッグ用のプリント文を追加

		// (1) エグジット処理：現在の signal.SignalDate に対応するエグジット日があるか確認
		for exitDate, exits := range exitMap {
			if signal.SignalDate.After(exitDate) {
				for _, exit := range exits {

					// 総資金に利益率 / 100% × ポジションサイズ × エントリー価格を加算
					profitInAmount := exit.ProfitLoss / 100 * exit.PositionSize * exit.EntryPrice
					portfolioValue += int(profitInAmount)

					// 1トレードあたりの損益率の単純加算
					if exit.ProfitLoss > 0 {
						winCount++
					}
					totalCount++
					tradeResults = append(tradeResults, exit) // トレード結果を保存
					delete(activeTrades, exit.Symbol)         // ホールド解除

					// 資金を更新した後の状態を表示
					fmt.Printf("%s (%s) 銘柄:%-4s [エントリ:%5.0f - %5.0f :エグジット] 損益/トレード: %4.1f%%, 総資産:%10d\n",
						exit.ExitDate.Format("2006-01-02"),
						exit.EntryDate.Format("2006-01-02"),
						exit.Symbol,
						exit.EntryPrice,
						exit.ExitPrice,
						exit.ProfitLoss,
						portfolioValue)

				}
				delete(exitMap, exitDate) // エグジット済みのデータを削除
			}
		}

		// (2) 既にホールド中ならスキップ
		if _, holding := activeTrades[signal.Symbol]; holding {
			// fmt.Println("ホールド中:", signal)
			continue
		}
		// 現在のポジションを差し引いた使用可能資金を計算
		availableFunds = portfolioValue
		for _, trade := range activeTrades {
			positionValue := trade.EntryPrice * trade.PositionSize
			availableFunds -= int(positionValue)
		}

		// (3) シンボルのデータを検索してエントリー処理
		for _, symbolData := range response.SymbolData {
			if symbolData.Symbol != signal.Symbol {
				// fmt.Printf("スキップ: 銘柄 %s は既にホールド中\n", signal.Symbol) // 【デバッグ用】 ホールド中のためスキップをログ出力
				continue
			}
			// ---- エントリー資金計算 ----
			positionSize, entryPrice, entryCost, err := determinePositionSize(portfolioValue, availableFunds, &symbolData.DailyData, signal.SignalDate)
			if err != nil || entryCost == 0 {
				// fmt.Println("エントリーコスト 0 のためスキップ") // 【デバッグ用】 エントリーコスト0でスキップをログ出力
				continue
			}

			// 使用可能資金を引く前にチェック
			availableFundsAfterTrade := availableFunds - int(entryCost)
			if availableFundsAfterTrade < 0 {
				// fmt.Println("使用可能資金不足のためシグナルをスキップ") // 資金不足でスキップをログ出力
				continue
			}
			availableFunds = availableFundsAfterTrade // 使用可能資金を引く

			// ---- トレード実行 ----
			// fmt.Println("トレード実行")
			purchaseDate, exitDate, profitLoss, _, exitPrice, err := singleTradingStrategy(
				&symbolData.DailyData, signal.SignalDate, stopLossPercentage, trailingStopTrigger, trailingStopUpdate,
			)
			// fmt.Println("exitDate:", exitDate)
			if err != nil {
				// fmt.Println("トレード実行 skip")
				continue
			}
			// ---- エントリー情報の保存 ----
			activeTrades[signal.Symbol] = tradeRecord{
				Symbol:       signal.Symbol,
				EntryDate:    purchaseDate,
				ExitDate:     exitDate,
				ProfitLoss:   profitLoss,
				EntryCost:    entryCost,
				PositionSize: positionSize,
				EntryPrice:   entryPrice,
				ExitPrice:    exitPrice,
			}
			// エグジット情報も `exitMap` に追加
			exitMap[exitDate] = append(exitMap[exitDate], tradeRecord{
				Symbol:       signal.Symbol,
				EntryDate:    purchaseDate,
				ExitDate:     exitDate,
				ProfitLoss:   profitLoss,
				EntryCost:    entryCost,
				PositionSize: positionSize,
				EntryPrice:   entryPrice,
				ExitPrice:    exitPrice,
			})
			// fmt.Println("exitMap[exitDate]:", exitMap)
		}
	}

	// 勝率の計算
	winRate := 0.0
	if totalCount > 0 {
		winRate = float64(winCount) / float64(totalCount) * 100
	}

	// 平均利益、平均損失の計算
	averageProfit, averageLoss := calculateAverages(tradeResults)
	// 最大ドローダウンの計算
	maxDrawdown := calculateMaxDrawdown(tradeResults)
	// シャープレシオの計算（リスク対リターンの指標）
	sharpeRatio := calculateSharpeRatio(tradeResults, 0)
	// リスク報酬比率の計算
	riskRewardRatio := 0.0
	if averageLoss != 0 {
		riskRewardRatio = averageProfit / math.Abs(averageLoss)
	}
	// 期待値の計算（トレード1回あたりの平均利益）
	expectedValue := 0.0
	if totalCount > 0 {
		expectedValue = ((winRate * averageProfit) - ((100 - winRate) * averageLoss)) / 100
	}

	// 最大連続利益と最大連続損失の計算
	maxConsecutiveProfit, maxConsecutiveLoss := calculateMaxConsecutive(tradeResults)

	return totalProfitLoss, winRate, maxConsecutiveProfit, maxConsecutiveLoss, winCount, totalCount - winCount, averageProfit, averageLoss, maxDrawdown, sharpeRatio, riskRewardRatio, expectedValue, nil
}
