package main

import (
	"fmt"
	"log"
	"os"

	"go-optimal-stop/experiment_proto" // パスをそのまま維持

	"google.golang.org/protobuf/proto"
)

// go run experiment_proto/print_response/main.go

func main() {
	// 確認したいファイル
	filePath := "data/ml_stock_response/latest_response.bin"

	// 表示する行数を指定する変数
	displayRows := 5

	// ファイルを読み込む
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// プロトコルバッファのメッセージをデシリアライズ
	var response experiment_proto.MLStockResponse
	if err := proto.Unmarshal(data, &response); err != nil {
		log.Fatalf("Failed to parse proto file: %v", err)
	}

	// レスポンスの内容を表示
	printResponse(&response, displayRows)
}

func printResponse(response *experiment_proto.MLStockResponse, displayRows int) {
	for _, symbolData := range response.SymbolData {
		fmt.Printf("Symbol: %s\n", symbolData.Symbol)
		fmt.Println("Daily Data:")
		printFirstAndLastN(dailyDataList(symbolData.DailyData), displayRows)

		fmt.Println("Signals:")
		for _, signal := range symbolData.Signals {
			fmt.Println(signal)
		}

		fmt.Println("Model Predictions:")
		for model, predictions := range symbolData.ModelPredictions {
			fmt.Printf("Model: %s\n", model)
			printFirstAndLastN(predictionsList(predictions.PredictionDates), displayRows)
		}
		fmt.Println()
	}
}

func printFirstAndLastN(data []string, n int) {
	if len(data) <= 2*n {
		for _, item := range data {
			fmt.Println(item)
		}
	} else {
		for _, item := range data[:n] {
			fmt.Println(item)
		}
		fmt.Println("...")
		for _, item := range data[len(data)-n:] {
			fmt.Println(item)
		}
	}
}

func dailyDataList(dailyData []*experiment_proto.MLDailyData) []string {
	var list []string
	for _, data := range dailyData {
		list = append(list, fmt.Sprintf("Date: %s, Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f, Volume: %d",
			data.Date, data.Open, data.High, data.Low, data.Close, data.Volume))
	}
	return list
}

func predictionsList(predictionDates []string) []string {
	var list []string
	list = append(list, predictionDates...)
	return list
}
