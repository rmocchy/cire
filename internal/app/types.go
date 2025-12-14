package app

// StructAnalysisResult は構造体の解析結果を保持する
type StructAnalysisResult struct {
	StructName  string                // 構造体名
	PackagePath string                // パッケージパス
	Fields      []FieldAnalysisResult // フィールドの解析結果
	Skipped     bool                  // 解析がスキップされたかどうか
	SkipReason  string                // スキップされた理由
}

// FieldAnalysisResult はフィールドの解析結果を保持する
type FieldAnalysisResult struct {
	Name                string                // フィールド名
	TypeName            string                // 型名
	PackagePath         string                // パッケージパス
	IsPointer           bool                  // ポインタ型かどうか
	IsInterface         bool                  // インターフェース型かどうか
	ResolvedStruct      *StructAnalysisResult // インターフェースから解決された構造体（再帰的）
	InterfaceSkipped    bool                  // インターフェースの解析がスキップされたか
	InterfaceSkipReason string                // インターフェースのスキップ理由
}
