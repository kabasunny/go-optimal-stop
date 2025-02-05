package trading

import (
	"errors"
	"time"

	"go-optimal-stop/internal/ml_stockdata"
)

// singleTradingStrategy: トレーディング戦略の実行
func singleTradingStrategy(data *[]ml_stockdata.InMLDailyData, startDate time.Time, stopLossPercentage, trailingStopTrigger, trailingStopUpdate float64) (time.Time, time.Time, float64, float64, float64, error) {
	d := *data

	// データが空の場合、エラーを返す
	if len(d) == 0 {
		return time.Time{}, time.Time{}, 0, 0, 0, errors.New("データが空です")
	}

	// データの最大日付を取得
	maxDate, err := parseDate(d[len(d)-1].Date)
	if err != nil {
		return time.Time{}, time.Time{}, 0, 0, 0, err
	}

	// 開始日がデータ範囲内になるまで、1日ずつ進める
	for !dateInData(d, startDate) {
		if startDate.After(maxDate) {
			return time.Time{}, time.Time{}, 0, 0, 0, errors.New("開始日がデータの範囲外です。無限ループを防ぐため、処理を中断")
		}
		startDate = startDate.AddDate(0, 0, 1)
	}

	// 購入日と購入価格を見つける
	purchaseDate, purchasePrice, err := findPurchaseDate(d, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, 0, 0, 0, err
	}

	// 終了日と終了価格を見つける
	endDate, endPrice, err := findExitDate(d, purchaseDate, purchasePrice, stopLossPercentage, trailingStopTrigger, trailingStopUpdate)
	if err != nil {
		return time.Time{}, time.Time{}, 0, 0, 0, err
	}

	// 利益率を計算
	profitLoss := roundDown((endPrice - purchasePrice) / purchasePrice * 100)

	// 購入日、終了日、利益率、購入価格、終了価格を返す
	return purchaseDate, endDate, profitLoss, purchasePrice, endPrice, nil
}
