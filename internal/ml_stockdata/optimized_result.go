package ml_stockdata

// OptimizedResult 構造体の定義
type OptimizedResult struct {
	StopLossPercentage   float64 // 初回はロスカット、以降はストップオーダーの設定値となる
	TrailingStopTrigger  float64 // トレイリングストップのトリガー閾値
	TrailingStopUpdate   float64 // トレイリングストップトリガー発動時の設定値
	ProfitLoss           float64 // 損益率
	WinRate              float64 // 勝率
	MaxConsecutiveProfit float64 // 連続して加算された最大利益の幅
	MaxConsecutiveLoss   float64 // 連続して加算された最大損失の幅
	PurchaseDate         string  // 購入日
	ExitDate             string  // 売却日
	ConsecutiveWins      int     // 連続益の回数
	ConsecutiveLosses    int     // 連続損の回数
	TotalWins            int     // 勝ち総回数
	TotalLosses          int     // 負け総回数
	AverageProfit        float64 // 平均利益率
	AverageLoss          float64 // 平均損失率
	MaxDrawdown          float64 // 最大ドローダウン
	SharpeRatio          float64 // シャープレシオ
	RiskRewardRatio      float64 // リスクリワード比
	ExpectedValue        float64 // 期待値
}
