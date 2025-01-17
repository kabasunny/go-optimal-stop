package main

import (
	"go-optimal-stop/experiment_proto"
)

func main() {
	filePath := "data/ml_stock_response/2025-01-17_11-12-14.bin"
	experiment_proto.RunOptimization(filePath)
}
