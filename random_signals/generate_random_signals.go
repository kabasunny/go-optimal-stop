package random_signals

import (
	"go-optimal-stop/internal/ml_stockdata"
	"math/rand"
	"sort"
	"time"
)

// 日付データからランダムにシグナルを選ぶ関数
func generateRandomSignals(data []ml_stockdata.InMLDailyData, numSignals int, seed ...int64) []string {

	var r *rand.Rand
	// seedが提供されている場合、ランダム数生成器にseedを使用
	if len(seed) > 0 {
		r = rand.New(rand.NewSource(seed[0]))
	} else {
		// seedが提供されていない場合、現在のUnixNano時間を使用
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	// データをシャッフルするためのコピーを作成
	shuffledData := make([]ml_stockdata.InMLDailyData, len(data))
	copy(shuffledData, data)
	// データをランダムにシャッフル
	r.Shuffle(len(shuffledData), func(i, j int) {
		shuffledData[i], shuffledData[j] = shuffledData[j], shuffledData[i]
	})

	// ランダムシグナルを格納するスライスを初期化
	signals := make([]string, 0, numSignals)
	// numSignalsの数だけランダムな日付を選択
	for i := 0; i < numSignals && i < len(shuffledData); i++ {
		signals = append(signals, shuffledData[i].Date)
	}

	// 日付順にシグナルを並べ替える
	sort.Strings(signals)

	return signals
}
