package trading

import (
	// 【デバッグ用】 fmt パッケージをインポート
	"go-optimal-stop/internal/ml_stockdata"
	"math"
	"time"
)

// determinePositionSize は、ATRに基づきポジションサイズとエントリー価格、エントリーコストを決定
func determinePositionSize(portfolioValue int, availableFundsInt int, entryPrice float64, commissionRate *float64, dailyData *[]ml_stockdata.InMLDailyData, signalDate time.Time) (float64, float64, error) {

	const unitSize = 100 // 単元数

	availableFunds := float64(availableFundsInt)

	// 引数に変更
	// エントリー価格を取得
	// _, entryPrice, err := findPurchaseDate(*dailyData, signalDate)
	// if err != nil {
	// 	return 0, 0, 0, err
	// }

	// ATRを計算
	atr := calculateATR(dailyData, signalDate)

	// リスク許容度を定義（例: 使用可能な資金の2%）
	riskPerTrade := 0.01 * availableFunds

	// ポジションサイズを計算（リスク許容度 / ATR）
	positionSize := 0.0
	if atr != 0 {
		positionSize = riskPerTrade / (atr * 2)
	}

	// ポジションサイズを調整して単元数の倍数にする
	positionSize = math.Floor(positionSize/float64(unitSize)) * float64(unitSize)

	// 手数料を加味してエントリーコストを計算
	entryCost := entryPrice * positionSize
	commission := entryCost * (*commissionRate / 100)
	totalEntryCost := entryCost + commission

	// 使用可能な資金に対してエントリーコストが足りるか確認
	if totalEntryCost <= availableFunds && totalEntryCost <= float64(portfolioValue)/4 {
		// 条件を満たす場合、ポジションサイズ、エントリー価格、エントリーコストを返す
		return positionSize, totalEntryCost, nil
	} else {
		// 条件を満たさない場合はエントリーしない
		return 0, 0, nil
	}
}

// calculateATR は、過去一定期間のATR（Average True Range）を計算する
func calculateATR(dailyData *[]ml_stockdata.InMLDailyData, signalDate time.Time) float64 {
	// ATRの計算ロジック（過去n日間のTrue Range平均）
	n := 14 // 計算に使用する日数
	trueRanges := make([]float64, 0, n)

	// fmt.Println("calculateATR 開始")                                              // 【デバッグ用】 関数開始をログ出力
	// fmt.Printf("  シグナル日付: %s, 計算期間: %d日\n", signalDate.Format("2006-01-02"), n) // 【デバッグ用】 シグナル日付と計算期間をログ出力

	// signalDate以前のn日間のデータを収集
	for i := len(*dailyData) - 1; i >= 0; i-- {
		if len(trueRanges) >= n {
			break
		}

		data := (*dailyData)[i]
		date, _ := time.Parse("2006-01-02", data.Date)
		if date.After(signalDate) || date.Equal(signalDate) {
			continue
		}

		if i == 0 {
			break
		}

		yesterdayData := (*dailyData)[i-1]
		trueRange := calculateTrueRange(data, yesterdayData)
		trueRanges = append(trueRanges, trueRange)
		// fmt.Printf("    日付: %s, True Range: %.2f, trueRanges (len: %d): %v\n", data.Date, trueRange, len(trueRanges), trueRanges) // 【デバッグ用】 各日のTrue Rangeと trueRanges の状態をログ出力
	}
	// fmt.Printf("  収集した True Ranges: %v\n", trueRanges) // 【デバッグ用】 収集した True Ranges をログ出力

	if len(trueRanges) == 0 {
		// fmt.Println("  警告: True Ranges が空のため、ATR をデフォルト値 1.5 に設定") // 【デバッグ用】 True Ranges が空の場合の警告ログ
		return 1.5 // データが不足している場合のデフォルト値
	}

	// ATRを計算
	sum := 0.0
	for _, tr := range trueRanges {
		sum += tr
	}
	atr := sum / float64(len(trueRanges))
	// fmt.Printf("calculateATR 完了: ATR: %.2f (合計: %.2f, 日数: %d)\n", atr, sum, len(trueRanges)) // 【デバッグ用】 ATR 計算結果をログ出力
	return atr
}

// calculateTrueRange は、前日と当日のデータに基づいてTrue Rangeを計算
func calculateTrueRange(today, yesterday ml_stockdata.InMLDailyData) float64 {
	highLow := today.High - today.Low
	highClose := math.Abs(today.High - yesterday.Close)
	lowClose := math.Abs(today.Low - yesterday.Close)
	trueRange := math.Max(highLow, math.Max(highClose, lowClose))
	// fmt.Printf("calculateTrueRange: High-Low: %.2f, |High-前日Close|: %.2f, |Low-前日Close|: %.2f, True Range: %.2f\n", highLow, highClose, lowClose, trueRange) // 【デバッグ用】 True Range 計算の詳細をログ出力
	return trueRange
}
