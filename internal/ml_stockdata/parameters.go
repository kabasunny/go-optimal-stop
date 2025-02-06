// internal/stockdata/parameters.go

package ml_stockdata

// Parameter 構造体の定義 単発用シミュ用
type Parameter struct {
	StopLossPercentage  float64
	TrailingStopTrigger float64
	TrailingStopUpdate  float64
}

// Parameters 構造体の定義 組合せ探索シミュ用
type Parameters struct {
	StopLossPercentages  []float64
	TrailingStopTriggers []float64
	TrailingStopUpdates  []float64
}

// SetStopLoss メソッド
func (p *Parameters) SetStopLoss(start, end, step float64) {
	p.StopLossPercentages = generateRange(start, end, step)
}

// SetTrailingStop メソッド
func (p *Parameters) SetTrailingStop(start, end, step float64) {
	p.TrailingStopTriggers = generateRange(start, end, step)
}

// SetTrailingStopUpdate メソッド
func (p *Parameters) SetTrailingStopUpdate(start, end, step float64) {
	p.TrailingStopUpdates = generateRange(start, end, step)
}

// generateRange 関数: 範囲とステップから値を生成
func generateRange(start, end, step float64) []float64 {
	var rangeValues []float64
	for value := start; value <= end; value += step {
		rangeValues = append(rangeValues, value)
	}
	return rangeValues
}
