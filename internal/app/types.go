package app

// NodeType はノードの種類を表す
type NodeType int

const (
	NodeTypeStruct NodeType = iota
	NodeTypeInterface
)

// InitFunctionInfo は初期化関数の情報を保持する
type InitFunctionInfo struct {
	Name        string // 関数名
	PackagePath string // パッケージパス
}

// FieldNode はフィールドを表すインターフェース
type FieldNode interface {
	GetFieldName() string
	NodeType() NodeType
}

// StructNode は構造体ノードを表す（構造体の定義とフィールドを保持）
type StructNode struct {
	FieldName     string             // フィールド名（ルートの場合は空）
	StructName    string             // 構造体名
	PackagePath   string             // パッケージパス
	InitFunctions []InitFunctionInfo // 構造体を返す初期化関数
	Fields        []FieldNode        // フィールドのノード
	Skipped       bool               // 解析がスキップされたかどうか
	SkipReason    string             // スキップされた理由
}

func (s *StructNode) GetFieldName() string {
	return s.FieldName
}

func (s *StructNode) NodeType() NodeType {
	return NodeTypeStruct
}

// InterfaceNode はインターフェースフィールドを表す
type InterfaceNode struct {
	FieldName      string      // フィールド名
	TypeName       string      // インターフェース型名
	PackagePath    string      // パッケージパス
	ResolvedStruct *StructNode // 解決された構造体
	Skipped        bool        // 解決がスキップされたか
	SkipReason     string      // スキップされた理由
}

func (i *InterfaceNode) GetFieldName() string {
	return i.FieldName
}

func (i *InterfaceNode) NodeType() NodeType {
	return NodeTypeInterface
}
