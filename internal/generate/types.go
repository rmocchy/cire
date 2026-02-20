package generate

// WireData は wire.go テンプレートに渡すデータ
type WireData struct {
	PackageName  string
	Imports      []string
	ProviderSets []ProviderSetData
}

// ProviderSetData は各 Provider セットのデータ
type ProviderSetData struct {
	StructName string
	Providers  []string
}
