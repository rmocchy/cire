package packages

// FieldInfo は構造体のフィールド情報を保持する
type FieldInfo struct {
	Name        string // フィールド名
	TypeName    string // 型名（例: "UserService", "UserRepository"）
	PackagePath string // importに使ったパッケージパス（例: "github.com/rmocchy/convinient_wire/sample/basic/service"）
	IsPointer   bool   // ポインタ型かどうか
	IsInterface bool   // インターフェース型かどうか
}

// StructFieldsInfo は構造体とそのフィールド情報を保持する
type StructFieldsInfo struct {
	StructName string      // 構造体名
	Fields     []FieldInfo // フィールド情報のリスト
}

// FunctionInfo は関数情報を保持する
type FunctionInfo struct {
	Name        string // 関数名
	PackagePath string // パッケージパス
}
