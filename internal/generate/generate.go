package generate

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"slices"
)

// 生成に必要な型定義
type GenerateConfig struct {
	PackageName string
	StructSets  []StructSet
}

type StructSet struct {
	RootStructName string
	Providers      []Provider
}

type Provider struct {
	PkgPath string
	Name    string
}

func (c *GenerateConfig) AddStructSet(rootStructName string, providers []Provider) {
	c.StructSets = append(c.StructSets, StructSet{
		RootStructName: rootStructName,
		Providers:      providers,
	})
}

func (c *GenerateConfig) SetPackageName(pkgName string) {
	c.PackageName = pkgName
}

func (c *GenerateConfig) Generate() ([]byte, error) {
	imports := make(map[string]bool)
	for _, set := range c.StructSets {
		for _, provider := range set.Providers {
			imports[provider.PkgPath] = true
		}
	}

	importList := make([]string, 0, len(imports))
	for imp := range imports {
		importList = append(importList, imp)
	}

	providers := make(map[string][]string)
	for _, set := range c.StructSets {
		providerNames := make([]string, 0, len(set.Providers))
		for _, provider := range set.Providers {
			providerNames = append(providerNames, provider.Name)
		}
		providers[set.RootStructName] = providerNames
	}

	// providerをソート
	providerSet := make([]ProviderSetData, 0, len(providers))
	for key := range providers {
		providerNames := providers[key]
		slices.Sort(providerNames)
		providerSet = append(providerSet, ProviderSetData{
			StructName: key,
			Providers:  providerNames,
		})
	}

	data := WireData{
		PackageName:  c.PackageName,
		Imports:      importList,
		ProviderSets: providerSet,
	}

	tmpl := template.Must(template.New("wire").Parse(wireTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		formatted = buf.Bytes()
	}
	return formatted, nil
}
