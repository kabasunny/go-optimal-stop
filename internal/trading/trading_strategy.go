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
	sort.Slice(signals, func(i, j int) bool {
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
					// 資金を更新する前の状態を表示（必要に応じて）
					// fmt.Printf("エグジット前 - シンボル: %s, originalTotalFunds: %d, availableFunds: %d\n", exit.Symbol, originalTotalFunds, availableFunds)

					// 使用可能資金にExitPrice × ポジションサイズを加算
					// exitAmount := exit.ExitPrice * exit.PositionSize
					// availableFunds += int(exitAmount)

					// 総資金に利益率 / 100% × ポジションサイズ × エントリー価格を加算
					profitInAmount := exit.ProfitLoss / 100 * exit.PositionSize * exit.EntryPrice
					portfolioValue += int(profitInAmount)

					// デバッグ用の変数表示
					// fmt.Printf("デバッグ情報 - シンボル: %s\n", exit.Symbol)
					// fmt.Printf("  exitAmount: %.2f, exit.ExitPrice: %.2f, exit.PositionSize: %.2f\n", exitAmount, exit.ExitPrice, exit.PositionSize)
					// fmt.Printf("  profitInAmount: %.2f, exit.ProfitLoss: %.2f, exit.PositionSize: %.2f, exit.EntryPrice: %.2f\n", profitInAmount, exit.ProfitLoss, exit.PositionSize, exit.EntryPrice)

					// その他の更新
					totalProfitLoss += exit.ProfitLoss
					if exit.ProfitLoss > 0 {
						winCount++
					}
					totalCount++
					tradeResults = append(tradeResults, exit) // トレード結果を保存
					delete(activeTrades, exit.Symbol)         // ホールド解除

					// 資金を更新した後の状態を表示
					// 資金を更新した後の状態を表示
					fmt.Printf("エグジット日: %s エントリ日: %s シンボル: %s, ポートフォリオ: %d, エントリ金額: %.0f, 購入可能枠: %d\n", exit.ExitDate.Format("2006-01-02"), exit.EntryDate.Format("2006-01-02"), exit.Symbol, portfolioValue, exit.EntryCost, availableFunds)
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
			positionSize, entryPrice, entryCost, err := determinePositionSize(portfolioValue, &symbolData.DailyData, signal.SignalDate)
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
				ExitPrice:    exitPrice,
				PositionSize: positionSize,
				EntryPrice:   entryPrice, // EntryPriceを保存
			}
			// エグジット情報も `exitMap` に追加
			exitMap[exitDate] = append(exitMap[exitDate], tradeRecord{
				Symbol:       signal.Symbol,
				EntryDate:    purchaseDate,
				ExitDate:     exitDate,
				ProfitLoss:   profitLoss,
				EntryCost:    entryCost,
				ExitPrice:    exitPrice,
				PositionSize: positionSize,
				EntryPrice:   entryPrice, // EntryPriceを保存
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
