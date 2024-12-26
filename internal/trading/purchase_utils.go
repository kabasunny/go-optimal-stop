// internal/trading/purchase_utils.go

package trading

import (
	"errors"
	"go-optimal-stop/internal/ml_stockdata"
	"time"
)

// findPurchaseDate 関数: 購入日を見つける
func findPurchaseDate(data []ml_stockdata.Data, startDate time.Time) (time.Time, float64, error) {
	for _, day := range data {
		parsedDate, err := parseDate(day.Date)
		if err != nil {
			return time.Time{}, 0, err
		}
		if parsedDate.Equal(startDate) {
			return parsedDate, day.Open, nil
		}
	}
	return time.Time{}, 0, errors.New("購入日が見つかりませんでした")
}
