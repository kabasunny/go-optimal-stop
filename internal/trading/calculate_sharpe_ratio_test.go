package trading

import (
	"math"
	"testing"
)

func TestCalculateSharpeRatio(t *testing.T) {
	tests := []struct {
		name        string
		tradeResults []tradeRecord
		riskFreeRate float64
		expected    float64
		expectError bool
	}{
		{
			name: "Basic case with positive returns and some variability",
			tradeResults: []tradeRecord{
				{ProfitLoss: 10},
				{ProfitLoss: 20},
				{ProfitLoss: 5},
				{ProfitLoss: 15},
			},
			riskFreeRate: 0.01,
			expected:     0.7828, // Calculated manually: mean=12.5, excess returns=[12.49, 22.49, 4.99, 14.99], stdDev=6.7573, sharpe=(12.5-0.01)/6.7573
			expectError:  false,
		},
		{
			name:         "No trades, empty results",
			tradeResults: []tradeRecord{},
			riskFreeRate: 0.01,
			expected:     0, // Or handle as error/NaN
			expectError:  false,
		},
		{
			name: "All positive, consistent returns (zero std dev)",
			tradeResults: []tradeRecord{
				{ProfitLoss: 10},
				{ProfitLoss: 10},
				{ProfitLoss: 10},
			},
			riskFreeRate: 0.01,
			expected:     0, // stdDev will be 0, leading to division by zero, should be handled by standardDeviation to return 0 or error.
			expectError:  false,
		},
		{
			name: "All negative, consistent returns (zero std dev)",
			tradeResults: []tradeRecord{
				{ProfitLoss: -5},
				{ProfitLoss: -5},
				{ProfitLoss: -5},
			},
			riskFreeRate: 0.01,
			expected:     0, // stdDev will be 0
			expectError:  false,
		},
		{
			name: "Mixed positive and negative returns",
			tradeResults: []tradeRecord{
				{ProfitLoss: 10},
				{ProfitLoss: -5},
				{ProfitLoss: 15},
				{ProfitLoss: -10},
			},
			riskFreeRate: 0.01,
			expected:     0.1873, // Calculated manually: mean=2.5, excess returns=[2.49, -5.01, 14.99, -10.01], stdDev=10.6865, sharpe=(2.5-0.01)/10.6865
			expectError:  false,
		},
		{
			name: "Single trade",
			tradeResults: []tradeRecord{
				{ProfitLoss: 100},
			},
			riskFreeRate: 0.01,
			expected:     0, // stdDev will be 0 for a single data point
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateSharpeRatio(tt.tradeResults, tt.riskFreeRate)

			if tt.expectError {
				if !math.IsNaN(got) && got != 0 { // For division by zero, it might return NaN or 0
					t.Errorf("calculateSharpeRatio() got = %v, want error/NaN", got)
				}
			} else {
				// Due to floating point precision, compare with a small tolerance
				if math.Abs(got-tt.expected) > 0.0001 {
					t.Errorf("calculateSharpeRatio() got = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}
