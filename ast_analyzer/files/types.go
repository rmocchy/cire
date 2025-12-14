package file

// StructInfo は構造体の情報を保持する構造体
type StructInfo struct {
	Name      string // 構造体名
	IsPointer bool   // ポインタ型かどうか
}

// FunctionInfo は関数の情報を保持する構造体
type FunctionInfo struct {
	Name        string
	ReturnTypes []StructInfo // エラー型以外の返り値の構造体情報
}
