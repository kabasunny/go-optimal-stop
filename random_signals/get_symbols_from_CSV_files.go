package random_signals

import (
	"os"
	"path/filepath"
	"strings"
)

// GetSymbolsFromCSVFiles は、指定されたディレクトリ内のCSVファイル名からシンボルを取得します
func GetSymbolsFromCSVFiles(dir string) ([]string, error) {
	var symbols []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".csv" {
			// 拡張子を除いたファイル名を取得してシンボルとして使用
			symbol := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			symbols = append(symbols, symbol)
		}
	}

	return symbols, nil
}
