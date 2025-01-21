package random_signals

import (
	"encoding/csv"
	"go-optimal-stop/internal/ml_stockdata"
	"os"
	"strconv"
	"time"
)

// CSVファイルを読み込み、データをstockdata.Data構造体のスライスに変換
func loadCSV(filePath string, startDate string) ([]ml_stockdata.InMLDailyData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // 可変長の行を許可
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// 基準の日付をパース
	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	var data []ml_stockdata.InMLDailyData
	for i, record := range records {
		// ヘッダー行をスキップ
		if i == 0 {
			continue
		}

		date := record[1] // 日付がインデックス1にある

		// 日付をパース
		recordDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}

		// 基準の日付以降のデータのみを追加
		if recordDate.Before(startDateTime) {
			continue
		}

		open, err := strconv.ParseFloat(record[3], 64) // 修正: 開始価格はインデックス3
		if err != nil {
			return nil, err
		}
		high, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}
		low, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, err
		}
		close, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseInt(record[7], 10, 64)
		if err != nil {
			return nil, err
		}

		data = append(data, ml_stockdata.InMLDailyData{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}
	return data, nil
}
