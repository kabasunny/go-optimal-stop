package trading

import (
	"fmt"
	"go-optimal-stop/internal/ml_stockdata"
	"math"
	"sort"
	"time"
)

// TradingStrategy 関数は、与えられた株価データとトレーディングパラメータに基づいて最適なパラメータの組み合わせを見つける
func TradingStrategy(response *ml_stockdata.InMLStockResponse, totalFunds *int, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64, options ...bool) (ml_stockdata.OptimizedResult, error) {
	var result ml_stockdata.OptimizedResult
	var verbose bool

	// verbose オプションをチェック
	if len(options) > 0 {
		verbose = options[0]
	}

	// パラメータを保存
	result.StopLossPercentage = stopLossPercentage
	result.TrailingStopTrigger = trailingStopTrigger
	result.TrailingStopUpdate = trailingStopUpdate

	// エントリー可能金額までのエントリー順序を決定する
	signals := []struct {
		Symbol     string
		SignalDate time.Time
		Priority   int64
	}{}

	// 各銘柄のシグナルを取得し、日付順にソート
	for _, symbolData := range response.SymbolData {
		if len(symbolData.Signals) < 1 {
			continue
		}
		for _, signal := range symbolData.Signals {
			date, err := parseDate(signal)
			if err != nil {
				if verbose {
					fmt.Println("signal skip")
				}
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
	for _, signal := range signals {
		for exitDate, exits := range exitMap {
			if signal.SignalDate.After(exitDate) {
				for _, exit := range exits {
					profitInAmount := exit.ProfitLoss / 100 * exit.PositionSize * exit.EntryPrice
					portfolioValue += int(profitInAmount)
					if exit.ProfitLoss > 0 {
						result.TotalWins++
					} else {
						result.TotalLosses++
					}
					totalCount++
					tradeResults = append(tradeResults, exit)
					delete(activeTrades, exit.Symbol)
					// 最適パラメータ時だけ表示したいので、if追加
					if verbose {
						fmt.Printf("%s (%s) 銘柄:%-4s [エントリ:%5.0f - %5.0f :エグジット] 損益/トレード: %4.1f%%, 総資産:%10d\n",
							exit.ExitDate.Format("2006-01-02"),
							exit.EntryDate.Format("2006-01-02"),
							exit.Symbol,
							exit.EntryPrice,
							exit.ExitPrice,
							exit.ProfitLoss,
							portfolioValue)
					}
				}
				// マップから削除してリソースを解放
				delete(exitMap, exitDate)
			}
		}

		if _, holding := activeTrades[signal.Symbol]; holding {
			continue
		}

		availableFunds = portfolioValue
		for _, trade := range activeTrades {
			positionValue := trade.EntryPrice * trade.PositionSize
			availableFunds -= int(positionValue)
		}

		for _, symbolData := range response.SymbolData {
			if symbolData.Symbol != signal.Symbol {
				continue
			}
			positionSize, entryPrice, entryCost, err := determinePositionSize(portfolioValue, availableFunds, &symbolData.DailyData, signal.SignalDate)
			if err != nil || entryCost == 0 {
				continue
			}
			availableFundsAfterTrade := availableFunds - int(entryCost)
			if availableFundsAfterTrade < 0 {
				continue
			}
			availableFunds = availableFundsAfterTrade

			purchaseDate, exitDate, profitLoss, _, exitPrice, err := singleTradingStrategy(
				&symbolData.DailyData, signal.SignalDate, stopLossPercentage, trailingStopTrigger, trailingStopUpdate,
			)
			if err != nil {
				continue
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
			}

			activeTrades[signal.Symbol] = record
			exitMap[exitDate] = append(exitMap[exitDate], record)
		}
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
