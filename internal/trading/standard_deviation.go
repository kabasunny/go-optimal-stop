package trading

import "math"

// standardDeviation 関数は、与えられたデータセットの標準偏差を計算
func standardDeviation(data []float64) float64 {
	mean := 0.0
	// データセットの平均値を計算
	for _, value := range data {
		mean += value
	}
	mean /= float64(len(data))

	// データセットの分散を計算
	variance := 0.0
	for _, value := range data {
		variance += (value - mean) * (value - mean)
	}
	variance /= float64(len(data))

	// 標準偏差を計算し、結果を返す
	return math.Sqrt(variance)
}
