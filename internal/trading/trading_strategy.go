package trading

import (
	"fmt"
	"go-optimal-stop/internal/ml_stockdata"
	"math"
	"sort"
	"time"
)

// TradingStrategy 関数は、与えられた株価データとトレーディングパラメータに基づいて総利益、勝率、その他の指標を返す
func TradingStrategy(response *ml_stockdata.InMLStockResponse, totalFunds *int, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (float64, float64, float64, float64, int, int, float64, float64, float64, float64, float64, float64, error) {
	signals := []struct {
		Symbol     string
		SignalDate time.Time
		Priority   int64
	}{}

	// 各銘柄のシグナルを取得し、日付順にソート
	for _, symbolData := range response.SymbolData {
		for _, signal := range symbolData.Signals {
			date, err := parseDate(signal)
			if err != nil {
				fmt.Println("signal skip")
				continue // 日付の解析に失敗した場合はスキップ
			}
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

	activeTrades := make(map[string]tradeRecord) // 各シンボルのホールド状態
	originalTotalFunds := *totalFunds            // 総資金の初期化（コピーを作成）
	availableFunds := *totalFunds                // 使用可能な資金の初期化
	totalProfitLoss := 0.0                       // 全体の利益を追跡
	winCount, totalCount := 0, 0                 // 勝ちトレード数と総トレード数
	var tradeResults []tradeRecord               // トレード結果を保持するスライス

	// シンボルごとのエグジット情報を保持するマップ
	exitMap := make(map[time.Time][]tradeRecord)

	// ---- 既存のエグジット情報を exitMap に格納 ----
	for _, record := range activeTrades {
		exitMap[record.ExitDate] = append(exitMap[record.ExitDate], record)
	}
	fmt.Printf("len(signals): %d\n", len(signals))
	// ---- シグナルの処理 ----
	for _, signal := range signals {
		fmt.Println("シグナル処理中:", signal) // デバッグ用のプリント文を追加

		// (1) エグジット処理：現在の signal.SignalDate に対応するエグジット日があるか確認
		for exitDate, exits := range exitMap {
			if signal.SignalDate.After(exitDate) {
				for _, exit := range exits {
					fmt.Printf("エグジット処理中: シンボル: %s, エグジット日付: %v, 利益損失: %v\n", exit.Symbol, exit.ExitDate, exit.ProfitLoss) // デバッグ用のプリント文を追加
					availableFunds += int(exit.ExitPrice)                                                                  // 使用可能資金を戻す
					originalTotalFunds += int(exit.ProfitLoss)                                                             // 総資金を逐次更新
					totalProfitLoss += exit.ProfitLoss
					if exit.ProfitLoss > 0 {
						winCount++
					}
					totalCount++
					tradeResults = append(tradeResults, exit) // トレード結果を保存
					delete(activeTrades, exit.Symbol)         // ホールド解除
					fmt.Println("ホールド解除:", exit.Symbol)       // デバッグ用のプリント文を追加
				}
				delete(exitMap, exitDate) // エグジット済みのデータを削除
			}
		}

		// (2) 既にホールド中ならスキップ
		if _, holding := activeTrades[signal.Symbol]; holding {
			fmt.Println("ホールド中:", signal)
			continue
		}

		// (3) シンボルのデータを検索してエントリー処理
		for _, symbolData := range response.SymbolData {
			if symbolData.Symbol != signal.Symbol {
				fmt.Printf("スキップ: 銘柄 %s は既にホールド中\n", signal.Symbol) // 【デバッグ用】 ホールド中のためスキップをログ出力
				continue
			}
			// ---- エントリー資金計算 ----
			_, _, entryCost, err := determinePositionSize(originalTotalFunds, &symbolData.DailyData, signal.SignalDate)
			if err != nil || entryCost == 0 {
				fmt.Println("エントリーコスト 0 のためスキップ") // 【デバッグ用】 エントリーコスト0でスキップをログ出力
				continue
			}
			availableFunds -= int(entryCost) // 使用可能資金を引く

			// ---- トレード実行 ----
			fmt.Println("トレード実行")
			purchaseDate, exitDate, profitLoss, _, exitPrice, err := singleTradingStrategy(
				&symbolData.DailyData, signal.SignalDate, stopLossPercentage, trailingStopTrigger, trailingStopUpdate,
			)
			fmt.Println("exitDate:", exitDate)
			if err != nil {
				fmt.Println("トレード実行 skip")
				continue
			}
			// ---- エントリー情報の保存 ----
			activeTrades[signal.Symbol] = tradeRecord{
				Symbol:     signal.Symbol,
				EntryDate:  purchaseDate,
				ExitDate:   exitDate,
				ProfitLoss: profitLoss,
				EntryCost:  entryCost,
				ExitPrice:  exitPrice,
			}
			// エグジット情報も `exitMap` に追加
			exitMap[exitDate] = append(exitMap[exitDate], tradeRecord{
				Symbol:     signal.Symbol,
				EntryDate:  purchaseDate,
				ExitDate:   exitDate,
				ProfitLoss: profitLoss,
				EntryCost:  entryCost,
				ExitPrice:  exitPrice,
			})
			fmt.Println("exitMap[exitDate]:", exitMap)
		}
	}

	// 計算処理（勝率・リスク管理指標）
	winRate := 0.0
	if totalCount > 0 {
		winRate = float64(winCount) / float64(totalCount) * 100
	}

	averageProfit, averageLoss := calculateAverages(tradeResults)
	maxDrawdown := calculateMaxDrawdown(tradeResults)
	sharpeRatio := calculateSharpeRatio(tradeResults, 0)
	riskRewardRatio := 0.0
	if averageLoss != 0 {
		riskRewardRatio = averageProfit / math.Abs(averageLoss)
	}
	expectedValue := 0.0
	if totalCount > 0 {
		expectedValue = ((winRate * averageProfit) - ((100 - winRate) * averageLoss)) / 100
	}

	maxConsecutiveProfit, maxConsecutiveLos := calculateMaxConsecutive(tradeResults)

	return totalProfitLoss, winRate, maxConsecutiveProfit, maxConsecutiveLos, winCount, totalCount - winCount, averageProfit, averageLoss, maxDrawdown, sharpeRatio, riskRewardRatio, expectedValue, nil
}
