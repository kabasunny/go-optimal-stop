// cmd/utils.go

package main

import (
	"math/rand"
	"time"

	"go-optimal-stop/internal/stockdata"
)

// 日付データからランダムに10個のシグナルを選ぶ関数
func generateRandomSignals(data []stockdata.Data, numSignals int, seed ...int64) []string {
	var r *rand.Rand
	if len(seed) > 0 {
		r = rand.New(rand.NewSource(seed[0]))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	shuffledData := make([]stockdata.Data, len(data))
	copy(shuffledData, data)
	r.Shuffle(len(shuffledData), func(i, j int) {
		shuffledData[i], shuffledData[j] = shuffledData[j], shuffledData[i]
	})
	signals := make([]string, 0, numSignals)
	for i := 0; i < numSignals && i < len(shuffledData); i++ {
		signals = append(signals, shuffledData[i].Date)
	}
	return signals
}
