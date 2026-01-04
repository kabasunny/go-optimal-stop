# ポジションサイジングと意思決定ロジックの解説

このドキュメントでは、提供されたGoコードベースから抽出した、トレーディングにおけるポジションサイジング（建玉量の決定）と意思決定（エントリーおよびエグジット）に関するロジ-ックについて解説します。

## 概要

このトレーディングシステムは、事前に外部から与えられた売買シグナル（買いの日付）に基づいて動作します。システムのコアロジックは、シグナルを生成することではなく、「どのように取引を実行し、リスクを管理するか」という点にあります。

主な意思決定のフローは以下の通りです。

1.  **シグナルの受信**: 取引すべき銘柄と日付のリストを受け取ります。
2.  **ポジションサイズの決定**: ポートフォリオ全体のリスクを考慮し、ATR（Average True Range）を用いて1取引あたりの適切な建玉量を計算します。
3.  **エントリー**: シグナルで指定された日の始値で株式を購入します。
4.  **エグジット**: あらかじめ定義された2種類の方法（静的なストップロス、または利益が出た後に有効になるトレーリングストップ）に基づいて、保有ポジションを売却します。

以下に、このロジックを実装している主要な関数とコードを解説します。

---

## 1. ポジションサイジングのロジック (`determine_position_size.go`)

ポジションサイズは、有名なリスク管理手法である「Van Tharpモデル」に基づいて決定されます。これは、1回のトレードで許容できる最大損失額を事前に決め、それに基づいて建玉量を調整するアプローチです。

このロジックは `DeterminePositionSize` 関数に実装されています。

### コード

```go
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
```

### 解説

1.  **ATRの計算 (`calculateATR`)**:
    *   `calculateATR` は、過去14日間のATR（Average True Range）を計算します。ATRは株価のボラティリティ（変動幅）を測る指標です。
    *   `calculateTrueRange` で1日あたりの真の変動幅を計算し、その平均をATRとしています。

2.  **1トレードあたりの許容損失額の計算**:
    *   `allowedLoss := float64(portfolioValue) * (param.RiskPercentage / 100)`
    *   ポートフォリオ全体の価値（`portfolioValue`）に対して、パラメータで指定された一定の割合（`RiskPercentage`）を掛け合わせ、この1回の取引で失ってもよい最大金額を算出します。例えば、ポートフォリオが100万円で`RiskPercentage`が2%なら、許容損失額は2万円となります。

3.  **1株あたりのリスク額（ストップロス幅）**:
    *   `stopLossAmount := atr * param.ATRMultiplier`
    *   ATRにパラメータで指定された倍率（`ATRMultiplier`）を掛けて、1株あたりの想定リスク（エントリー価格から損切りラインまでの幅）を決定します。ATRが大きい（変動が激しい）銘柄ほど、リスク額は大きくなります。

4.  **初期ポジションサイズの計算**:
    *   `initialPositionSize := allowedLoss / stopLossAmount`
    *   「1トレードあたりの許容損失額」を「1株あたりのリスク額」で割ることで、購入すべき株数（ポジションサイズ）を計算します。

5.  **ポジションサイズの調整**:
    *   `initialPositionSize = math.Floor(initialPositionSize/float64(unitSize)) * float64(unitSize)`
    *   計算されたポジションサイズを、取引所の最小取引単位（`unitSize`、ここでは100株）の倍数になるように切り捨てます。
    *   さらに、算出したポジションの合計コスト（`initialTotalEntryCost`）が、利用可能な現金（`availableFunds`）やポートフォリオのリスク許容範囲（`riskLimitEntryCost`）を超えないように、ポジションサイズを下方修正します。これにより、資金不足や過大なリスクを避けます。

最終的に、すべての制約を満たした `positionSize`（株数）と `totalEntryCost`（合計コスト）が返されます。

---

## 2. エグジット（売却）の意思決定ロジック (`stop_order_utils.go`)

ポジションをいつ手仕舞うかは、`findExitDate` 関数によって決定されます。この関数は、静的なストップロスと、利益が出た後に発動するトレーリングストップの2つのルールを実装しています。

### コード

```go
package trading

import (
	"go-optimal-stop/internal/ml_stockdata"
	"time"
)

// roundUp 関数: 四捨五入（切り上げ）
func roundUp(value float64) float64 {
	return float64(int(value*10+1)) / 10
}

// roundDown 関数: 四捨五入（切り捨て）
func roundDown(value float64) float64 {
	return float64(int(value*10)) / 10
}

// findExitDate: 売却日と売却価格を決定
func findExitDate(data []ml_stockdata.InMLDailyData, purchaseDate time.Time, purchasePrice float64, param *ml_stockdata.Parameter) (time.Time, float64, error) {
	var endDate time.Time
	var endPrice float64

	// パラメータを保存
	stopLossPercentage := param.StopLossPercentage
	trailingStopTrigger := param.TrailingStopTrigger
	trailingStopUpdate := param.TrailingStopUpdate

	// ストップロスとトレーリングストップの閾値を計算
	stopLossThreshold := roundDown(purchasePrice * (1 - stopLossPercentage/100))
	trailingStopTriggerPrice := roundUp(purchasePrice * (1 + trailingStopTrigger/100))

	// purchaseDate 以降のデータを取得（スライスを最適化）
	var startIndex int
	for i, day := range data {
		parsedDate, err := parseDate(day.Date)
		if err != nil {
			return time.Time{}, 0, err
		}
		if !parsedDate.Before(purchaseDate) {
			startIndex = i
			break
		}
	}
	filteredData := data[startIndex:]

	// トレーリングストップの監視
	for _, day := range filteredData {
		parsedDate, err := parseDate(day.Date)
		if err != nil {
			return time.Time{}, 0, err
		}

		openPrice := day.Open
		lowPrice := day.Low
		closePrice := day.Close

		// ストップロスに到達した場合
		if lowPrice <= stopLossThreshold || openPrice <= stopLossThreshold {
			endPrice = stopLossThreshold
			endDate = parsedDate
			break
		}

		// トレーリングストップのトリガーをチェック
		if closePrice >= trailingStopTriggerPrice {
			trailingStopTriggerPrice = roundUp(closePrice * (1 + trailingStopTrigger/100))
			stopLossThreshold = roundDown(closePrice * (1 - trailingStopUpdate/100))
		}
	}

	// 途中で売却しなかった場合、最終データを採用
	if endDate.IsZero() {
		lastIndex := len(filteredData) - 1
		endPrice = filteredData[lastIndex].Close
		endDate, _ = parseDate(filteredData[lastIndex].Date)
	}

	return endDate, endPrice, nil
}
```

### 解説

1.  **初期損切りラインの設定**:
    *   `stopLossThreshold := roundDown(purchasePrice * (1 - stopLossPercentage/100))`
    *   購入価格（`purchasePrice`）から、パラメータで指定された静的な損切り率（`stopLossPercentage`）だけ低い価格を、最初の損切りラインとして設定します。

2.  **トレーリングストップの発動条件**:
    *   `trailingStopTriggerPrice := roundUp(purchasePrice * (1 + trailingStopTrigger/100))`
    *   購入価格が一定の利益率（`trailingStopTrigger`）に達したら、トレーリングストップが発動します。この価格が `trailingStopTriggerPrice` です。

3.  **日々の価格監視**:
    *   購入日以降、日足データのループ処理に入ります。
    *   `if lowPrice <= stopLossThreshold || openPrice <= stopLossThreshold`:
        *   その日の安値または始値が現在の損切りライン (`stopLossThreshold`) を下回った場合、その時点で損切りが確定し、ループを抜けます。
    *   `if closePrice >= trailingStopTriggerPrice`:
        *   終値がトレーリングストップの発動価格を上回った場合、トレーリングストップが有効になります。
        *   損切りライン (`stopLossThreshold`) が、現在の終値から一定割合（`trailingStopUpdate`）だけ低い価格に引き上げられます。
        *   同時に、次のトレーリングストップ発動価格もさらに高い価格に更新されます。
    *   これにより、利益が伸びる限り損切りラインを切り上げ続け、損失を限定しつつ利益を確保しようとします。

4.  **最終的な売却**:
    *   もし期間中に一度も損切りラインに掛からなかった場合は、データ期間の最終日の終値で売却されます。

---

## 3. 全体戦略の実行 (`trading_strategy.go`)

`TradingStrategy` 関数は、これまでのロジックを統合し、バックテスト全体を管理するオーケストレーターです。

### コード（抜粋）

```go
// TradingStrategy 関数は、与えられた株価データとトレーディングパラメータに基づいて最適なパラメータの組み合わせを見つける
func TradingStrategy(response *ml_stockdata.InMLStockResponse, totalFunds *int, param *ml_stockdata.Parameter, commissionRate *float64, options ...bool) (ml_stockdata.OptimizedResult, error) {
    // ... (変数の初期化) ...

	// シグナルを日付順、優先順にソート
	sort.Slice(signals, func(i, j int) bool {
		// ...
	})

    // ... (ポートフォリオ変数の初期化) ...

	// ---- シグナルの処理 ----
	for _, signal := range signals {
        // ... (既存ポジションの決済処理) ...

        // ... (利用可能資金の計算) ...

		for _, symbolData := range response.SymbolData {
			if symbolData.Symbol != signal.Symbol {
				continue
			}
            // singleTradingStrategyは内部でfindPurchaseDateとfindExitDateを呼び出す
			purchaseDate, exitDate, profitLoss, entryPrice, exitPrice, err := singleTradingStrategy(
				&symbolData.DailyData, signal.SignalDate, param,
			)
			if err != nil {
				continue
			}

            // ポジションサイズを決定
			positionSize, entryCost, err := DeterminePositionSize(param, portfolioValue, availableFunds, entryPrice, commissionRate, &symbolData.DailyData, signal.SignalDate)
			if err != nil || entryCost == 0 {
				continue
			}

            // ... (取引記録の作成とアクティブトレードへの追加) ...
		}
	}

    // ... (最終的なパフォーマンス指標の計算) ...

	return result, nil
}
```

### 解説

1.  **シグナルのソート**:
    *   受け取ったすべてのシグナルを、日付と優先度に基づいてソートします。これにより、時系列に沿ってバックテストが実行されます。

2.  **シグナルのループ処理**:
    *   ソートされたシグナルを一つずつ処理します。
    *   ループの冒頭で、現在のシグナル日付より前に決済されるべき既存のポジションを処理し、ポートフォリオの価値（`portfolioValue`）と利用可能資金（`availableFunds`）を更新します。

3.  **ポジションサイズ決定とエントリー**:
    *   `singleTradingStrategy` を呼び出して、エントリー価格や仮の exitDate を取得します。
    *   次に、`DeterminePositionSize` を呼び出して、計算されたエントリー価格と現在のポートフォリオ状況に基づき、実際に購入する株数を決定します。
    *   資金が不足している場合や、計算されたポジションサイズが0の場合は、取引を見送ります。

4.  **結果の記録**:
    *   取引が実行された場合、その内容（エントリー日、コスト、株数など）を `activeTrades` マップに記録します。
    *   ループが完了した後、すべての取引結果を集計し、勝率、平均利益/損失、最大ドローダウンなどのパフォーマンス指標を計算して返します。
