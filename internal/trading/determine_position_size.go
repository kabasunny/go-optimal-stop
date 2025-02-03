package trading

import (
	"fmt" // 【デバッグ用】 fmt パッケージをインポート
	"go-optimal-stop/internal/ml_stockdata"
	"math"
	"time"
)

// determinePositionSize は、ATRに基づきポジションサイズとエントリー価格、エントリーコストを決定
func determinePositionSize(portfolioValue int, dailyData *[]ml_stockdata.InMLDailyData, signalDate time.Time) (float64, float64, float64, error) {
	const commissionRate = 0.2 // 手数料率（例: 0.1%）
	const unitSize = 100       // 単元数

	// fmt.Println("determinePositionSize 開始")                                               // 【デバッグ用】 関数開始をログ出力
	// fmt.Printf("  総資金: %d, シグナル日付: %s\n", currentFunds, signalDate.Format("2006-01-02")) // 【デバッグ用】 初期資金とシグナル日付をログ出力

	// エントリー価格を取得
	_, entryPrice, err := findPurchaseDate(*dailyData, signalDate) // entryData も取得するように修正
	if err != nil {
		fmt.Printf("  findPurchaseDate エラー: %v\n", err) // 【デバッグ用】 findPurchaseDate エラーをログ出力
		return 0, 0, 0, err                             // エントリー価格の取得に失敗した場合はエラーを返す
	}
	// fmt.Printf("  findPurchaseDate 完了: エントリー価格: %.2f, 日付: %s\n", entryPrice, entryData.Date) // 【デバッグ用】 エントリー価格と日付をログ出力

	// ATRを計算
	atr := calculateATR(dailyData, signalDate)
	// fmt.Printf("  calculateATR 完了: 2ATR: %.2f\n", atr*2) // 【デバッグ用】 ATR をログ出力

	// リスク許容度を定義（例: 総資金の2%）
	riskPerTrade := 0.01 * float64(portfolioValue)
	// fmt.Printf("  リスク許容額: %.2f (総資金の2%%)\n", riskPerTrade) // 【デバッグ用】 リスク許容額をログ出力

	// ポジションサイズを計算（リスク許容度 / ATR）
	positionSize := 0.0
	if atr != 0 { // ATR が 0 でないことを確認
		positionSize = riskPerTrade / (atr * 2)
	} else {
		fmt.Println("  警告: ATR が 0 のため、ポジションサイズを 0 に設定") // 【デバッグ用】 ATR が 0 の場合の警告ログ
	}
	// fmt.Printf("  ポジションサイズ (調整前): %.2f\n", positionSize) // 【デバッグ用】 調整前のポジションサイズをログ出力

	// ポジションサイズを調整して単元数の倍数にする
	positionSize = math.Floor(positionSize/float64(unitSize)) * float64(unitSize) // Floor を使用して単元未満を切り捨て
	// fmt.Printf("  ポジションサイズ (調整後): %.2f (単元数: %d)\n", positionSize, unitSize)      // 【デバッグ用】 調整後のポジションサイズと単元数をログ出力

	// 手数料を加味してエントリーコストを計算
	entryCost := entryPrice * positionSize
	commission := entryCost * (commissionRate / 100)
	totalEntryCost := entryCost + commission
	// fmt.Printf("  エントリーコスト計算: エントリー価格: %.2f, ポジションサイズ: %.2f, 手数料: %.2f, 合計: %.2f\n", entryPrice, positionSize, commission, totalEntryCost) // 【デバッグ用】 エントリーコスト計算の詳細をログ出力

	// 総資金に対してエントリーコストが足りなければエントリーコストは0にする
	if totalEntryCost > float64(portfolioValue) {
		// fmt.Println("  エントリーコストが初期資金を超えるため、エントリーコストを 0 に設定") // 【デバッグ用】 資金不足でエントリーコストが 0 になる場合のログ
		return 0, 0, 0, nil
	}
	if positionSize == 0 { // ポジションサイズが 0 の場合もエントリーコストを 0 にする
		fmt.Println("  ポジションサイズが 0 のため、エントリーコストを 0 に設定") // 【デバッグ用】 ポジションサイズが 0 の場合のエントリーコストを 0 にするログ
		return 0, 0, 0, nil
	}

	// fmt.Printf("determinePositionSize 完了: ポジションサイズ: %.2f, エントリー価格: %.2f, エントリーコスト: %.2f\n", positionSize, entryPrice, totalEntryCost) // 【デバッグ用】 関数終了時の結果をログ出力
	return positionSize, entryPrice, totalEntryCost, nil
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
