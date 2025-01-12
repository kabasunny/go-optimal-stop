package main

import (
	"fmt"
	"io"
	"log"
	"os"

	parquet "github.com/parquet-go/parquet-go"
)

// 実行コマンド　go run .\parquet\show_parquet.go
func main() {
	fileName := "parquet/1570_2025-01-12.parquet"

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open Parquet file: %v", err)
	}
	defer file.Close()

	// ファイルサイズを取得
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	size := fileInfo.Size()

	// Parquet ファイルを開く
	pf, err := parquet.OpenFile(file, size)
	if err != nil {
		log.Fatalf("Failed to open Parquet file as parquet.File: %v", err)
	}

	// スキーマを取得して表示
	schema := pf.Schema()
	fmt.Println("Parquet Schema:")
	for _, field := range schema.Fields() {
		fmt.Printf("Name: %s, Type: %s\n", field.Name(), field.Type())
	}

	// 指定した行数
	numRowsToRead := 20

	// 行を読み込む
	for _, rowGroup := range pf.RowGroups() {
		rows := rowGroup.Rows()
		fmt.Printf("rows %s\n", rows)
		defer rows.Close()

		// parquet.Rowのスライスを事前に確保
		row := make([]parquet.Row, numRowsToRead)

		for i := 0; i < numRowsToRead; {
			// 行を読み込む
			n, err := rows.ReadRows(row)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("Failed to read row: %v", err)
			}

			// 読み込んだ行を表示
			for j := 0; j < n; j++ {
				if i >= numRowsToRead {
					break
				}
				fmt.Printf("Row %d: %v\n", i+1, row[j])
				i++
			}
		}
	}
}
