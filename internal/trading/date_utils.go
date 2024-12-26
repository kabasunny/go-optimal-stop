// internal/trading/date_utils.go

package trading

import (
	"go-optimal-stop/internal/ml_stockdata"
	"time"
)

// parseDate 関数: 日付文字列を time.Time 型に変換
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// dateInData 関数: 特定の日付がデータに存在するか確認
func dateInData(data []ml_stockdata.Data, date time.Time) bool {
	for _, day := range data {
		parsedDate, err := parseDate(day.Date)
		if err != nil {
			return false
		}
		if parsedDate.Equal(date) {
			return true
		}
	}
	return false
}
