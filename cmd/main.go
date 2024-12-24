// cmd/main.go

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"go-optimal-stop/internal/optimization"
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
		low, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
		close, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}

		data = append(data, stockdata.Data{
			Date:  date,
			Open:  open,
			Low:   low,
			Close: close,
		})
	}
	return data, nil
}

func main() {
	filePath := "7203.T_stock_data.csv"

	// CSVファイルを読み込み
	data, err := loadCSV(filePath)
	if err != nil {
		fmt.Printf("CSVファイルの読み込みエラー: %v\n", err)
		return
	}

	tradeStartDate := "2023-01-01" // 開始日を文字列として指定

	// パラメータの最適化を実行
	bestResult, worstResult, _ := optimization.OptimizeParameters(&data, tradeStartDate)

	// 結果を表示
	fmt.Printf("最良の結果: %+v\n", bestResult)
	fmt.Printf("最悪の結果: %+v\n", worstResult)
}
