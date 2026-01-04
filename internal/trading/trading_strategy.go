package trading

import (
	"fmt"
	"go-optimal-stop/internal/ml_stockdata"
	"math"
	"sort"
	"strconv"
	"time"
)

// TradingStrategy 関数は、与えられた株価データとトレーディングパラメータに基づいて最適なパラメータの組み合わせを見つける
func TradingStrategy(response *ml_stockdata.InMLStockResponse, totalFunds *int, param *ml_stockdata.Parameter, commissionRate *float64, options ...bool) (ml_stockdata.OptimizedResult, error) {
	var result ml_stockdata.OptimizedResult
	var verbose bool

	// verbose オプションをチェック
	if len(options) > 0 {
		verbose = options[0]
	}

	// パラメータを保存
	result.StopLossPercentage = param.StopLossPercentage
	result.TrailingStopTrigger = param.TrailingStopTrigger
	result.TrailingStopUpdate = param.TrailingStopUpdate
	result.ATRMultiplier = param.ATRMultiplier
	result.RiskPercentage = param.RiskPercentage

	// エントリー可能金額までのエントリー順序を決定する
	signals := []struct {
		Symbol     string
		SignalDate time.Time
		Priority   int64
	}{}

	// 各銘柄のシグナルを取得し、日付順にソート（買いシグナルと売りシグナルの両方）
	for _, symbolData := range response.SymbolData {
		if len(symbolData.Signals) < 1 {
			if verbose {
				fmt.Printf("銘柄 %s: シグナルなし\n", symbolData.Symbol)
			}
			continue
		}
		if verbose {
			fmt.Printf("銘柄 %s: %d個のシグナルを処理中\n", symbolData.Symbol, len(symbolData.Signals))
		}
		for _, signal := range symbolData.Signals {
			date, err := parseDate(signal.Date)
			if err != nil {
				if verbose {
					fmt.Printf("シグナル日付解析エラー: %s\n", signal.Date)
				}
				continue // 日付の解析に失敗した場合はスキップ
			}

			// 買いシグナルと売りシグナルの両方を処理対象に追加
			signals = append(signals, struct {
				Symbol     string
				SignalDate time.Time
				Priority   int64
			}{symbolData.Symbol, date, symbolData.Priority})

			signalType := "買い"
			if signal.Type == ml_stockdata.SignalTypeSell {
				signalType = "売り"
			}
			if verbose {
				fmt.Printf("%sシグナル追加: %s (%s)\n", signalType, symbolData.Symbol, date.Format("2006-01-02"))
			}
		}
	}

	// シグナルを日付順、優先順にソート
	sort.Slice(signals, func(i, j int) bool {
		if signals[i].SignalDate.Equal(signals[j].SignalDate) {
			return signals[i].Priority < signals[j].Priority
		}
		return signals[i].SignalDate.Before(signals[j].SignalDate)
	})

	activeTrades := make(map[string]tradeRecord)
	originalTotalFunds := *totalFunds
	portfolioValue := originalTotalFunds
	availableFunds := portfolioValue
	totalCount := 0
	var tradeResults []tradeRecord

	// シンボルごとのエグジット情報を保持するマップ
	exitMap := make(map[time.Time][]tradeRecord)

	// ---- シグナルの処理 ----
	if verbose {
		fmt.Printf("処理対象シグナル数: %d\n", len(signals))
	}
	for _, signal := range signals {
		if verbose {
			fmt.Printf("シグナル処理中: %s (%s)\n", signal.Symbol, signal.SignalDate.Format("2006-01-02"))
		}

		// 現在のシグナルが買いシグナルか売りシグナルかを判定
		var currentSignalType ml_stockdata.SignalType = ml_stockdata.SignalTypeBuy // デフォルト
		for _, symbolData := range response.SymbolData {
			if symbolData.Symbol == signal.Symbol {
				for _, sig := range symbolData.Signals {
					if sig.Date == signal.SignalDate.Format("2006-01-02") {
						currentSignalType = sig.Type
						break
					}
				}
				break
			}
		}

		if currentSignalType == ml_stockdata.SignalTypeSell {
			// 売りシグナル処理
			if verbose {
				fmt.Printf("売りシグナル処理: %s\n", signal.Symbol)
			}
			if trade, exists := activeTrades[signal.Symbol]; exists {
				if signal.SignalDate.After(trade.EntryDate) {
					if verbose {
						fmt.Printf("売りシグナル実行: %s\n", signal.Symbol)
					}
					// 売却価格を取得
					exitPrice, err := findPriceAtDate(response, signal.Symbol, signal.SignalDate)
					if err != nil {
						if verbose {
							fmt.Printf("売却価格取得エラー: %v\n", err)
						}
						continue
					}

					// 利益/損失を計算
					profitLoss := (exitPrice - trade.EntryPrice) / trade.EntryPrice * 100
					profitInAmount := profitLoss / 100 * trade.PositionSize * trade.EntryPrice

					// 勝ち負けのカウント、勝ちの場合は税金20パーセント
					if profitLoss > 0 {
						result.TotalWins++
						profitInAmount *= 0.8
					} else {
						result.TotalLosses++
					}

					portfolioValue += int(profitInAmount)
					totalCount++

					// 取引記録を更新
					trade.ExitDate = signal.SignalDate
					trade.ExitPrice = exitPrice
					trade.ProfitLoss = profitLoss
					trade.ExitReason = "SellSignal"
					trade.SignalBased = true

					tradeResults = append(tradeResults, trade)
					delete(activeTrades, signal.Symbol)

					if verbose {
						entryCost := trade.EntryPrice * trade.PositionSize * (1 + *commissionRate/100)
						exitValue := trade.ExitPrice * trade.PositionSize * (1 - *commissionRate/100)
						profitPerPortfolio := (exitValue - entryCost) / float64(trade.PortfolioValue) * 100
						positionRate := entryCost / float64(trade.PortfolioValue) * 100

						fmt.Printf(" [%-4s](%10s) %10s : %9.0f - %9.0f (%5.0f)[ %11.0f (%6.1f%%) - %11.0f ] %8.1f%%, %8.2f%%, %10s [SELL SIGNAL]\n",
							trade.Symbol,
							trade.EntryDate.Format("2006-01-02"),
							trade.ExitDate.Format("2006-01-02"),
							trade.EntryPrice,
							trade.ExitPrice,
							trade.PositionSize,
							entryCost,
							positionRate,
							exitValue,
							trade.ProfitLoss,
							profitPerPortfolio,
							formatNumber(portfolioValue))
					}
				}
			}
			continue
		}

		// 買いシグナル処理（既存のロジック）
		// 売りシグナルのチェック（現在の日付で実行可能な売りシグナルを処理）
		// この部分は上記の統合処理で処理されるため削除

		for exitDate, exits := range exitMap {
			if signal.SignalDate.After(exitDate) {
				for _, exit := range exits {
					profitInAmount := exit.ProfitLoss / 100 * exit.PositionSize * exit.EntryPrice

					// 勝ち負けのカウント、勝ちの場合は税金20パーセント
					if exit.ProfitLoss > 0 {
						result.TotalWins++
						profitInAmount *= 0.8
					} else {
						result.TotalLosses++
					}

					portfolioValue += int(profitInAmount)
					totalCount++
					tradeResults = append(tradeResults, exit)
					delete(activeTrades, exit.Symbol)

					// 最適パラメータ時だけ表示したいので、if追加
					if verbose {
						entryCost := exit.EntryPrice * exit.PositionSize * (1 + *commissionRate/100)
						exitValue := exit.ExitPrice * exit.PositionSize * (1 - *commissionRate/100)
						profitPerPortfolio := (exitValue - entryCost) / float64(exit.PortfolioValue) * 100
						positionRate := entryCost / float64(exit.PortfolioValue) * 100

						fmt.Printf(" [%-4s](%10s) %10s : %9.0f - %9.0f (%5.0f)[ %11.0f (%6.1f%%) - %11.0f ] %8.1f%%, %8.2f%%, %10s\n",
							exit.Symbol,
							exit.EntryDate.Format("2006-01-02"),
							exit.ExitDate.Format("2006-01-02"),
							exit.EntryPrice,
							exit.ExitPrice,
							exit.PositionSize,
							entryCost,
							positionRate,
							exitValue,
							exit.ProfitLoss,    // 株価に対する割合
							profitPerPortfolio, // 総資産に対する割合
							formatNumber(portfolioValue))
					}
				}
				// マップから削除してリソースを解放
				delete(exitMap, exitDate)
			}
		}

		if _, holding := activeTrades[signal.Symbol]; holding {
			if verbose {
				fmt.Printf("既に保有中のためスキップ: %s\n", signal.Symbol)
			}
			continue
		}

		if verbose {
			fmt.Printf("新規エントリー検討: %s\n", signal.Symbol)
		}

		availableFunds = portfolioValue
		for _, trade := range activeTrades {
			positionValue := trade.EntryPrice * trade.PositionSize * (1 + *commissionRate)
			availableFunds -= int(positionValue)
		}

		for _, symbolData := range response.SymbolData {
			if symbolData.Symbol != signal.Symbol {
				continue
			}
			if verbose {
				fmt.Printf("銘柄データ見つかりました: %s\n", symbolData.Symbol)
			}
			purchaseDate, exitDate, profitLoss, entryPrice, exitPrice, err := singleTradingStrategy(
				&symbolData.DailyData, signal.SignalDate, param,
			)
			if err != nil {
				if verbose {
					fmt.Printf("singleTradingStrategy エラー: %v\n", err)
				}
				continue
			}
			if verbose {
				fmt.Printf("取引戦略結果: エントリー=%s, エグジット=%s, 利益率=%.2f%%, エントリー価格=%.2f\n",
					purchaseDate.Format("2006-01-02"), exitDate.Format("2006-01-02"), profitLoss, entryPrice)
			}

			positionSize, entryCost, err := DeterminePositionSize(param, portfolioValue, availableFunds, entryPrice, commissionRate, &symbolData.DailyData, signal.SignalDate)
			if err != nil || entryCost == 0 {
				if verbose {
					fmt.Printf("DeterminePositionSize 結果: positionSize=%.2f, entryCost=%.2f, err=%v\n", positionSize, entryCost, err)
				}
				continue
			}
			if verbose {
				fmt.Printf("ポジションサイズ決定成功: %.0f株, エントリーコスト: %.0f円\n", positionSize, entryCost)
			}
			availableFundsAfterTrade := availableFunds - int(entryCost)
			if verbose {
				fmt.Printf("資金チェック: 利用可能資金=%d円, エントリーコスト=%.0f円, 取引後資金=%d円\n",
					availableFunds, entryCost, availableFundsAfterTrade)
			}
			if availableFundsAfterTrade < 0 {
				if verbose {
					fmt.Println("利用可能資金不足のためエントリーをスキップ")
				}
				continue
			}

			// エントリー可能な資金の更新
			availableFunds = availableFundsAfterTrade

			if verbose {
				fmt.Printf("取引記録作成中: %s\n", signal.Symbol)
			}

			record := tradeRecord{
				Symbol:         signal.Symbol,
				EntryDate:      purchaseDate,
				ExitDate:       exitDate,
				ProfitLoss:     profitLoss,
				EntryCost:      entryCost,
				PositionSize:   positionSize,
				EntryPrice:     entryPrice,
				ExitPrice:      exitPrice,
				PortfolioValue: portfolioValue,
				ExitReason:     "AutoExit", // 自動売却（ストップロス/トレーリングストップ）
				SignalBased:    false,      // シグナルベースではない
			}

			if verbose {
				fmt.Printf("アクティブ取引に追加: %s (エントリー: %s, エグジット: %s)\n",
					record.Symbol, record.EntryDate.Format("2006-01-02"), record.ExitDate.Format("2006-01-02"))
			}

			activeTrades[signal.Symbol] = record
			exitMap[exitDate] = append(exitMap[exitDate], record)
		}
	}
	if verbose {
		fmt.Printf("トレード実行回数: %d\n", totalCount)
	}

	// 勝率の計算
	if totalCount > 0 {
		result.WinRate = float64(result.TotalWins) / float64(totalCount) * 100
	}

	// 平均利益、平均損失の計算
	result.AverageProfit, result.AverageLoss = calculateAverages(tradeResults)
	// 最大ドローダウンの計算
	result.MaxDrawdown, _ = calculateDrawdownAndDrawup(tradeResults)
	// シャープレシオの計算（リスク対リターンの指標）
	result.SharpeRatio = calculateSharpeRatio(tradeResults, 0)
	// リスク報酬比率の計算
	if result.AverageLoss != 0 {
		result.RiskRewardRatio = result.AverageProfit / math.Abs(result.AverageLoss)
	}
	// 期待値の計算（トレード1回あたりの平均利益）
	if totalCount > 0 {
		result.ExpectedValue = ((result.WinRate * result.AverageProfit) - ((100 - result.WinRate) * result.AverageLoss)) / 100
	}

	// 最大連続利益と最大連続損失の計算 現状使わない
	result.MaxConsecutiveProfit, result.MaxConsecutiveLoss = calculateMaxConsecutive(tradeResults)

	// 総利益の計算
	result.ProfitLoss = float64((portfolioValue - originalTotalFunds) * 100 / originalTotalFunds)

	return result, nil
}

func formatNumber(n int) string {
	in := strconv.FormatInt(int64(n), 10)
	out := make([]byte, len(in)+(len(in)-1)/3)
	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k > 2 {
			k = 0
			j--
			out[j] = ','
		}
	}
}

// findPriceAtDate 指定された日付の株価を取得する
func findPriceAtDate(response *ml_stockdata.InMLStockResponse, symbol string, targetDate time.Time) (float64, error) {
	for _, symbolData := range response.SymbolData {
		if symbolData.Symbol == symbol {
			for _, dailyData := range symbolData.DailyData {
				date, err := parseDate(dailyData.Date)
				if err != nil {
					continue
				}
				if date.Equal(targetDate) {
					return dailyData.Open, nil // 始値を使用
				}
			}
		}
	}
	return 0, fmt.Errorf("指定日付 %s の %s のデータが見つかりません", targetDate.Format("2006-01-02"), symbol)
}
