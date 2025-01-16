package main

import (
	"go-optimal-stop/experiment_proto"
)

func main() {
	filePath := "data/ml_stock_response/7203_2025-01-15.bin"
	experiment_proto.RunOptimization(filePath)
}
