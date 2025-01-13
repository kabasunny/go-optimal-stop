// internal/stockdata/stock_request.go

package ml_stockdata

// InMLStockRequest 構造体の定義
type InMLStockRequest struct {
	Symbols   []string // 銘柄コードのリスト
	StartDate string   // 開始日
	EndDate   string   // 終了日
}
