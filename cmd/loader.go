// cmd/loader.go

package main

import (
	"encoding/csv"
	"os"
	"strconv"

	"go-optimal-stop/internal/stockdata"
)

// CSVファイルを読み込み、データをstockdata.Data構造体のスライスに変換
func loadCSV(filePath string) ([]stockdata.Data, error) {
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

	var data []stockdata.Data
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

		data = append(data, stockdata.Data{
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
