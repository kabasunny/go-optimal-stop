package optimization

import (
	"go-optimal-stop/experiment_proto"
)

// extractModelNames 関数は、MLStockResponse プロトコルバッファのレスポンスからモデル名のリストを抽出
func extractModelNames(protoResponse *experiment_proto.MLStockResponse) []string {
	modelNameMap := make(map[string]struct{}) // 重複を避けるために、モデル名をキーとするマップを作成

	// protoResponse 内の各シンボルデータをループ処理
	for _, symbolData := range protoResponse.GetSymbolData() {
		// 各シンボルデータ内のモデル予測をループ処理して、モデル名をマップに追加
		for modelName := range symbolData.ModelPredictions {
			modelNameMap[modelName] = struct{}{} // モデル名をキーとしてマップに追加
		}
	}

	var modelNames []string // モデル名を格納するスライス
	// マップ内のモデル名をスライスに追加
	for modelName := range modelNameMap {
		modelNames = append(modelNames, modelName)
	}

	return modelNames // 抽出したモデル名のスライスを返す
}
