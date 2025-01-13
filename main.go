package main

import (
	"go-optimal-stop/experiment_proto"
)

func main() {
	filePath := "data/ml_stock_response/1570_2025-01-12.bin"
	experiment_proto.RunOptimization(filePath)
}
