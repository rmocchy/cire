package analyze

import (
	"go/types"

	"github.com/rmocchy/convinient_wire/internal/core"
)

// NodeType はノードの種類を表す
type NodeType int

const (
	NodeTypeStruct NodeType = iota
	NodeTypeInterface
	NodeTypeBuiltin
)

// FieldNode はフィールドを表すインターフェース
type FieldNode interface {
	GetFieldName() string
	NodeType() NodeType
}

// StructAnalyzer は構造体の解析を行うインターフェース
type StructAnalyzer interface {
	// AnalyzeStruct は構造体を解析する（エントリーポイント）
	AnalyzeStruct(structName string, packagePath core.PackagePath) (*StructNode, error)
	// AnalyzeNamedStructType は名前付き構造体型を解析する
	AnalyzeNamedStructType(structName string, packagePath core.PackagePath, structType *types.Struct) (*StructNode, error)
}

// StructNode は構造体ノードを表す（構造体の定義とフィールドを保持）
type StructNode struct {
	FieldName     string        // フィールド名（ルートの場合は空）
	StructName    string        // 構造体名
	PackagePath   string        // パッケージパス
	InitFunctions []*types.Func // 構造体を返す初期化関数
	Dependencies  []FieldNode   // 依存関係のフィールドノード
	Fields        []FieldNode   // フィールドのノード
	Skipped       bool          // 解析がスキップされたかどうか
	SkipReason    string        // スキップされた理由
}

func (s *StructNode) GetFieldName() string {
	return s.FieldName
}

func (s *StructNode) NodeType() NodeType {
	return NodeTypeStruct
}

// InterfaceNode はインターフェースフィールドを表す
type InterfaceNode struct {
	FieldName     string        // フィールド名
	TypeName      string        // インターフェース型名
	PackagePath   string        // パッケージパス
	InitFunctions []*types.Func // インターフェースを返す初期化関数
	Dependencies  []FieldNode   // 依存関係のフィールドノード
	Skipped       bool          // 解決がスキップされたか
	SkipReason    string        // スキップされた理由
}

func (i *InterfaceNode) GetFieldName() string {
	return i.FieldName
}

func (i *InterfaceNode) NodeType() NodeType {
	return NodeTypeInterface
}

type BuiltinNode struct {
	FieldName string // フィールド名
	TypeName  string // ビルトイン型名
}

func (b *BuiltinNode) GetFieldName() string {
	return b.FieldName
}

func (b *BuiltinNode) NodeType() NodeType {
	// ビルトイン型もNodeTypeStructとして扱う（必要に応じてNodeTypeを追加しても良い）
	return NodeTypeBuiltin
}
