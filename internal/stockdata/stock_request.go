// internal/stockdata/stock_request.go

package stockdata

// StockRequest 構造体の定義
type StockRequest struct {
	Symbols   []string // 銘柄コードのリスト
	StartDate string   // 開始日
	EndDate   string   // 終了日
}
