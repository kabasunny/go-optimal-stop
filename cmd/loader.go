// cmd/loader.go

package main

import (
	"encoding/csv"
	"os"
	"strconv"

	"go-optimal-stop/internal/ml_stockdata"
)

// CSVファイルを読み込み、データをstockdata.Data構造体のスライスに変換
func loadCSV(filePath string) ([]ml_stockdata.MLDailyData, error) {
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

	var data []ml_stockdata.MLDailyData
	for i, record := range records {
		// ヘッダー行をスキップ
		if i == 0 {
			continue
		}

		date := record[0]
		open, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		high, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		low, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
		close, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}
		adjClose, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseInt(record[6], 10, 64)
		if err != nil {
			return nil, err
		}

		data = append(data, ml_stockdata.MLDailyData{
			Date:     date,
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close,
			AdjClose: adjClose,
			Volume:   volume,
		})
	}
	return data, nil
}
