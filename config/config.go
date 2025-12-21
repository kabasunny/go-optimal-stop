package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Range defines the min, max, and step for a parameter.
type Range struct {
	Min  float64 `yaml:"min"`
	Max  float64 `yaml:"max"`
	Step float64 `yaml:"step"`
}

// OptimizationParameters holds the search space for all parameters.
type OptimizationParameters struct {
	StopLoss           Range `yaml:"stop_loss"`
	TrailingStop       Range `yaml:"trailing_stop"`
	TrailingStopUpdate Range `yaml:"trailing_stop_update"`
	ATRMultipliers     Range `yaml:"atr_multipliers"`
	RiskPercentages    Range `yaml:"risk_percentages"`
}

// Simulation holds the general settings for the simulation.
type Simulation struct {
	DataFilePath   string  `yaml:"data_file_path"`
	TotalFunds     int     `yaml:"total_funds"`
	CommissionRate float64 `yaml:"commission_rate"`
}

// Config is the top-level configuration struct.
type Config struct {
	Simulation            Simulation             `yaml:"simulation"`
	OptimizationParameters OptimizationParameters `yaml:"optimization_parameters"`
}

// LoadConfig reads the configuration from a YAML file.
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
