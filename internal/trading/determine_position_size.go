package trading

import (
	"fmt"
	"go-optimal-stop/internal/ml_stockdata"
	"math"
	"time"
)

// determinePositionSize は、ATRに基づきポジションサイズとエントリー価格、エントリーコストを決定
func DeterminePositionSize(param *ml_stockdata.Parameter, portfolioValue int, availableFundsInt int, entryPrice float64, commissionRate *float64, dailyData *[]ml_stockdata.InMLDailyData, signalDate time.Time) (float64, float64, error) {

	const unitSize = 100 // 単元数
	availableFunds := float64(availableFundsInt)

	// ATRを計算
	atr := calculateATR(dailyData, signalDate)
	// fmt.Println("ATR:", atr)

	// 許容損失額を計算 (ポートフォリオ価値のストップロス割合)
	allowedLoss := float64(portfolioValue) * (param.RiskPercentage / 100)

	// ストップロス幅をATRの2倍に設定（過去の価格変動の2倍の幅でストップロスを設定）
	stopLossAmount := atr * param.ATRMultiplier

	// 初期ポジションサイズを計算
	initialPositionSize := allowedLoss / stopLossAmount
	// fmt.Println("positionSize before unit size:", positionSize)

	// ポジションサイズを調整して最小単元の倍数にする
	initialPositionSize = math.Floor(initialPositionSize/float64(unitSize)) * float64(unitSize)
	// fmt.Println("positionSize after unit size:", positionSize)

	// 手数料を加味してエントリーコストを計算
	initialEntryCost := entryPrice * initialPositionSize
	commission := initialEntryCost * (*commissionRate / 100)
	initialTotalEntryCost := initialEntryCost + commission
	// fmt.Println("totalEntryCost:", totalEntryCost)

	// 使用可能な資金に対してエントリーコストが足りるか、
	// かつ、ポートフォリオのリスク許容範囲を超えないようにポジションサイズを調整
	maxPositionSize := initialPositionSize

	if initialTotalEntryCost > availableFunds {
		// 利用可能資金を超える場合、ポジションサイズを縮小
		maxPositionSize = math.Floor((availableFunds/(entryPrice*(1+(*commissionRate/100))))/float64(unitSize)) * float64(unitSize)
		// エントリーコストを再計算
		initialEntryCost = entryPrice * maxPositionSize
		commission = initialEntryCost * (*commissionRate / 100)
		initialTotalEntryCost = initialEntryCost + commission
	}

	// ポートフォリオのリスク許容範囲を超えないようにポジションサイズを調整
	riskLimitEntryCost := float64(portfolioValue) * param.RiskPercentage
	if initialTotalEntryCost > riskLimitEntryCost {
		maxPositionSize = math.Floor((riskLimitEntryCost/(entryPrice*(1+(*commissionRate/100))))/float64(unitSize)) * float64(unitSize)
		// エントリーコストを再計算
		initialEntryCost = entryPrice * maxPositionSize
		commission = initialEntryCost * (*commissionRate / 100)
		initialTotalEntryCost = initialEntryCost + commission
	}

	// 最終的なポジションサイズとエントリーコスト
	positionSize := maxPositionSize
	totalEntryCost := initialTotalEntryCost

	// ポジションサイズがゼロ以下の場合、エントリーしない
	if positionSize <= 0 {
		return 0, 0, nil
	}

	return positionSize, totalEntryCost, nil

}

// calculateATR は、過去一定期間のATR（Average True Range）を計算する
func calculateATR(dailyData *[]ml_stockdata.InMLDailyData, signalDate time.Time) float64 {
	// ATRの計算ロジック（過去n日間のTrue Range平均）
	n := 14 // 計算に使用する日数
	trueRanges := make([]float64, 0, n)

	// signalDate以前のn日間のデータを収集
	for i := len(*dailyData) - 1; i >= 1; i-- { // i >= 1 に変更 (yesterdayDataのために最低2つのデータが必要)
		data := (*dailyData)[i]
		date, _ := time.Parse("2006-01-02", data.Date)

		// signalDateより後のデータはスキップ
		if date.After(signalDate) {
			continue
		}

		// signalDate当日のデータもスキップ
		if date.Equal(signalDate) {
			continue
		}

		yesterdayData := (*dailyData)[i-1]
		trueRange := calculateTrueRange(data, yesterdayData)
		trueRanges = append([]float64{trueRange}, trueRanges...) // 先頭に追加
		//trueRanges = append(trueRanges, trueRange)
		if len(trueRanges) >= n {
			break
		}
	}

	if len(trueRanges) == 0 {
		fmt.Println("ATR計算に必要なデータが不足しています。エントリーを見送ります。")
		return 0 // ATRが計算できない場合は、0を返す（ポジションサイズが0になる）
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
