// internal/stockdata/stock_request.go

package ml_stockdata

// MLStockRequest 構造体の定義
type MLStockRequest struct {
	Symbols   []string // 銘柄コードのリスト
	StartDate string   // 開始日
	EndDate   string   // 終了日
}
