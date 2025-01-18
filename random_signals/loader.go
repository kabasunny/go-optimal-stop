package random_signals

import (
	"encoding/csv"
	"go-optimal-stop/internal/ml_stockdata"
	"os"
	"strconv"
)

// CSVファイルを読み込み、データをstockdata.Data構造体のスライスに変換
func loadCSV(filePath string) ([]ml_stockdata.InMLDailyData, error) {
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

	var data []ml_stockdata.InMLDailyData
	for i, record := range records {
		// ヘッダー行をスキップ
		if i == 0 {
			continue
		}

		date := record[1]                              // 日付がインデックス1にある
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

		// // デバッグ情報を追加
		// fmt.Printf("Date: %s, Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Volume: %d\n",
		// 	date, open, high, low, close, volume)

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
