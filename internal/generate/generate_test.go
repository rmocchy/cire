package generate

import (
	"go/format"
	"strings"
	"testing"
)

func TestGenerateConfig_Generate(t *testing.T) {
	tests := []struct {
		name        string
		config      *GenerateConfig
		wantErr     bool
		wantContain []string
	}{
		{
			name: "単一StructSet・単一Provider",
			config: &GenerateConfig{
				PackageName: "main",
				StructSets: []StructSet{
					{
						RootStructName: "App",
						Providers: []Provider{
							{PkgPath: "example.com/myapp/service", Name: "service.NewUserService"},
						},
					},
				},
			},
			wantErr: false,
			wantContain: []string{
				"package main",
				`"example.com/myapp/service"`,
				"var AppSet = wire.NewSet(",
				"service.NewUserService,",
				"wire.Struct(new(App)",
				"func InitializeApp()",
			},
		},
		{
			name: "複数StructSet",
			config: &GenerateConfig{
				PackageName: "mypkg",
				StructSets: []StructSet{
					{
						RootStructName: "Server",
						Providers: []Provider{
							{PkgPath: "example.com/handler", Name: "handler.NewHandler"},
						},
					},
					{
						RootStructName: "Client",
						Providers: []Provider{
							{PkgPath: "example.com/repo", Name: "repo.NewRepository"},
						},
					},
				},
			},
			wantErr: false,
			wantContain: []string{
				"package mypkg",
				"var ServerSet = wire.NewSet(",
				"handler.NewHandler,",
				"wire.Struct(new(Server)",
				"func InitializeServer()",
				"var ClientSet = wire.NewSet(",
				"repo.NewRepository,",
				"wire.Struct(new(Client)",
				"func InitializeClient()",
			},
		},
		{
			name: "複数Providerがソートされる",
			config: &GenerateConfig{
				PackageName: "main",
				StructSets: []StructSet{
					{
						RootStructName: "App",
						Providers: []Provider{
							{PkgPath: "example.com/svc", Name: "svc.NewZService"},
							{PkgPath: "example.com/svc", Name: "svc.NewAService"},
							{PkgPath: "example.com/svc", Name: "svc.NewMService"},
						},
					},
				},
			},
			wantErr: false,
			wantContain: []string{
				"svc.NewAService,",
				"svc.NewMService,",
				"svc.NewZService,",
			},
		},
		{
			name: "複数Providerの順序確認（AはZより前）",
			config: &GenerateConfig{
				PackageName: "main",
				StructSets: []StructSet{
					{
						RootStructName: "App",
						Providers: []Provider{
							{PkgPath: "example.com/svc", Name: "svc.NewZService"},
							{PkgPath: "example.com/svc", Name: "svc.NewAService"},
						},
					},
				},
			},
			wantErr: false,
			wantContain: []string{
				"svc.NewAService,",
				"svc.NewZService,",
			},
		},
		{
			name: "同一PkgPathの重複importが除外される",
			config: &GenerateConfig{
				PackageName: "main",
				StructSets: []StructSet{
					{
						RootStructName: "App",
						Providers: []Provider{
							{PkgPath: "example.com/svc", Name: "svc.NewFooService"},
							{PkgPath: "example.com/svc", Name: "svc.NewBarService"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "出力がgofmtされている",
			config: &GenerateConfig{
				PackageName: "main",
				StructSets: []StructSet{
					{
						RootStructName: "Foo",
						Providers: []Provider{
							{PkgPath: "example.com/foo", Name: "foo.NewFoo"},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.Generate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			output := string(got)

			for _, want := range tt.wantContain {
				if !strings.Contains(output, want) {
					t.Errorf("Generate() 出力に %q が含まれていません\n出力:\n%s", want, output)
				}
			}
		})
	}
}

func TestGenerateConfig_Generate_OutputIsFormatted(t *testing.T) {
	config := &GenerateConfig{
		PackageName: "main",
		StructSets: []StructSet{
			{
				RootStructName: "App",
				Providers: []Provider{
					{PkgPath: "example.com/svc", Name: "svc.NewService"},
				},
			},
		},
	}

	got, err := config.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	formatted, err := format.Source(got)
	if err != nil {
		t.Fatalf("生成コードがgofmtに通りません: %v", err)
	}

	if string(got) != string(formatted) {
		t.Errorf("Generate() の出力がgofmtされていません")
	}
}

func TestGenerateConfig_Generate_DuplicateImportsDeduped(t *testing.T) {
	config := &GenerateConfig{
		PackageName: "main",
		StructSets: []StructSet{
			{
				RootStructName: "App",
				Providers: []Provider{
					{PkgPath: "example.com/svc", Name: "svc.NewFoo"},
					{PkgPath: "example.com/svc", Name: "svc.NewBar"},
				},
			},
		},
	}

	got, err := config.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	output := string(got)
	count := strings.Count(output, `"example.com/svc"`)
	if count != 1 {
		t.Errorf("同一PkgPathのimportが重複しています: %d回出現", count)
	}
}
