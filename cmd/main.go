// cmd/main.go

package main

import (
	"fmt"
	"time"

	"go-optimal-stop/internal/data"
	"go-optimal-stop/internal/optimization"
)

func main() {
	// 例としてのデータとパラメータ
	data := []data.Data{
		{Date: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), Open: 100, Low: 98, Close: 102},
		{Date: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), Open: 102, Low: 100, Close: 104},
		{Date: time.Date(2023, 4, 4, 0, 0, 0, 0, time.UTC), Open: 112, Low: 110, Close: 114},
		{Date: time.Date(2023, 5, 5, 0, 0, 0, 0, time.UTC), Open: 102, Low: 100, Close: 104},
		// その他のデータ...
	}
	tradeStartDate := "2023-01-01" // 開始日を文字列として指定

	bestResult, worstResult, allResults := optimization.OptimizeParameters(&data, tradeStartDate)
	fmt.Printf("最良の結果: %+v\n", bestResult)
	fmt.Printf("最悪の結果: %+v\n", worstResult)
	fmt.Printf("全ての結果: %+v\n", allResults)
}
