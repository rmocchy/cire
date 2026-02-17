# Cire (Convenient Wire)

**Cire** は、Google Wire の依存性注入コードを自動生成する CLI ツールです。  
構造体の依存関係を解析し、`wire.go` ファイルを自動生成することで、Wire を使った DI の設定を簡単にします。

## 特徴

- 🔍 **自動依存関係解析**: 構造体のフィールドから依存関係ツリーを自動的に構築
- 📝 **wire.go 自動生成**: ProviderSet と Initialize 関数を含む `wire.go` を自動生成
- 🎯 **インターフェース対応**: インターフェース型のフィールドから適切なコンストラクタを検出
- 📊 **YAML 出力**: 依存関係の構造を YAML 形式で可視化（オプション）
- 🔧 **複数ルート構造体対応**: 1つのファイルで複数の DI ルート構造体を定義可能

## インストール

```bash
go install github.com/rmocchy/convinient_wire@latest
```

または、リポジトリをクローンしてビルド:

```bash
git clone https://github.com/rmocchy/convinient_wire.git
cd convinient_wire
make build
```

## クイックスタート

### 1. ルート構造体を定義

DI のルートとなる構造体を `cire.go` などのファイルに定義します。  
`//go:build wireinject` タグを付けることで、通常のビルドからは除外されます。

```go
//go:build wireinject
// +build wireinject

package main

import (
    "myapp/handler"
)

// App は依存関係のルート構造体
type App struct {
    Handler *handler.UserHandler
}
```

### 2. 各レイヤーの実装

通常の Go コードとして、各レイヤーの構造体とコンストラクタを実装します。

**Repository 層:**
```go
package repository

type UserRepository interface {
    FindByID(id int) (*User, error)
}

func NewUserRepository(config *Config) (UserRepository, error) {
    return &userRepositoryImpl{config: config}, nil
}
```

**Service 層:**
```go
package service

type UserService interface {
    GetUserInfo(id int) (string, error)
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userServiceImpl{repo: repo}
}
```

**Handler 層:**
```go
package handler

type UserHandler struct {
    service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
    return &UserHandler{service: service}
}
```

### 3. cire を実行

```bash
cire analyze --file ./cire.go
```

これにより `wire.go` が自動生成されます:

```go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
    "myapp/handler"
    "myapp/repository"
    "myapp/service"
)

var AppSet = wire.NewSet(
    repository.NewConfig,
    repository.NewUserRepository,
    service.NewUserService,
    handler.NewUserHandler,
    wire.Struct(new(App), "*"),
)

func InitializeApp() (*App, error) {
    wire.Build(AppSet)
    return nil, nil
}
```

### 4. Wire でコード生成

```bash
wire ./
```

`wire_gen.go` が生成され、アプリケーションで使用できるようになります。

## コマンド

### analyze

構造体の依存関係を解析し、`wire.go` を生成します。

```bash
cire analyze --file <ファイルパス> [--yaml]
```

| フラグ | 短縮形 | 説明 |
|--------|--------|------|
| `--file` | `-f` | 解析対象の Go ファイル（必須） |
| `--yaml` | `-y` | 依存関係を YAML ファイルに出力 |

**例:**
```bash
# wire.go のみ生成
cire analyze --file ./cire.go

# wire.go と cire.yaml を生成
cire analyze --file ./cire.go --yaml
```

## 依存関係の解析ルール

Cire は以下のルールで依存関係を解析します:

1. **構造体フィールド**: ポインタ型または値型の構造体フィールドを検出
2. **インターフェースフィールド**: インターフェース型のフィールドを検出し、対応するコンストラクタ (`New*` 関数) を探索
3. **コンストラクタ検出**: `New<型名>` という命名規則のコンストラクタを自動検出
4. **再帰解析**: 各依存の依存関係も再帰的に解析

## プロジェクト構成例

```
myapp/
├── cire.go          # ルート構造体定義
├── wire.go          # 自動生成される Wire 設定
├── wire_gen.go      # Wire が生成するファイル
├── main.go          # エントリーポイント
├── handler/
│   └── user_handler.go
├── service/
│   └── user_service.go
└── repository/
    └── user_repository.go
```

## サンプルプロジェクト

### Basic サンプル

シンプルな3層アーキテクチャの例:

```bash
make sample.basic
```

[sample/basic/](sample/basic/) を参照してください。

### Complex サンプル

複数のルート構造体と並列依存の例:

```bash
make sample.complex
```

[sample/complex/](sample/complex/) を参照してください。

## 開発

### ビルド

```bash
make build
```

### テスト

```bash
make test
```

### テストカバレッジ

```bash
make test.coverage
```

### クリーンアップ

```bash
make clean.all
```

## 要件

- Go 1.21 以上
- [Google Wire](https://github.com/google/wire)

## ライセンス

MIT License

## 作者

[@rmocchy](https://github.com/rmocchy)
