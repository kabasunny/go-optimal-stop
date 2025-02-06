package trading

import (
	"go-optimal-stop/internal/ml_stockdata"
	"math"
	"time"
)

// determinePositionSize は、ATRに基づきポジションサイズとエントリー価格、エントリーコストを決定
func DeterminePositionSize(StopLossPercentage float64, portfolioValue int, availableFundsInt int, entryPrice float64, commissionRate *float64, dailyData *[]ml_stockdata.InMLDailyData, signalDate time.Time) (float64, float64, error) {

	const unitSize = 100 // 単元数
	availableFunds := float64(availableFundsInt)

	// ATRを計算
	atr := calculateATR(dailyData, signalDate)

	// ポジションサイズを計算（ATR * 2 に基づく）
	positionSize := (atr * 2) / (StopLossPercentage / 100 * entryPrice)

	// ポジションサイズを調整して最小単元の倍数にする
	positionSize = math.Floor(positionSize/float64(unitSize)) * float64(unitSize)

	// 手数料を加味してエントリーコストを計算
	entryCost := entryPrice * positionSize
	commission := entryCost * (*commissionRate / 100)
	totalEntryCost := entryCost + commission

	// 使用可能な資金に対してエントリーコストが足りるか確認
	if totalEntryCost <= availableFunds && totalEntryCost <= float64(portfolioValue)/4 {
		// 条件を満たす場合、ポジションサイズ、エントリーコストを返す
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
	}

	if len(trueRanges) == 0 {
		return 1.5 // データが不足している場合のデフォルト値
	}

	// ATRを計算
	sum := 0.0
	for _, tr := range trueRanges {
		sum += tr
	}
	atr := sum / float64(len(trueRanges))
	return atr
}

// calculateTrueRange は、前日と当日のデータに基づいてTrue Rangeを計算
func calculateTrueRange(today, yesterday ml_stockdata.InMLDailyData) float64 {
	highLow := today.High - today.Low
	highClose := math.Abs(today.High - yesterday.Close)
	lowClose := math.Abs(today.Low - yesterday.Close)
	trueRange := math.Max(highLow, math.Max(highClose, lowClose))
	return trueRange
}
