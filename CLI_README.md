# Convenient Wire

Convenient Wireは、Go言語のプロジェクトにおける構造体の依存関係を解析し、YAML形式で出力するCLIツールです。

## 機能

- 指定された構造体の依存関係を再帰的に解析
- 初期化関数（コンストラクタ）の自動検出
- インターフェース、構造体、ビルトイン型のサポート
- YAML形式での依存関係の可視化

## インストール

```bash
go install github.com/rmocchy/convinient_wire@latest
```

または、リポジトリをクローンしてビルド:

```bash
git clone https://github.com/rmocchy/convinient_wire.git
cd convinient_wire
go build -o convinient_wire
```

## 使い方

### 基本的な使用方法

```bash
convinient_wire analyze --struct <構造体名> --package <パッケージパス>
```

### オプション

- `-s, --struct`: 解析する構造体名 (必須)
- `-p, --package`: 構造体が定義されているパッケージパス
- `-o, --output`: 出力ファイルパス (指定しない場合は標準出力)

### 例

#### 標準出力に結果を表示

```bash
convinient_wire analyze \
  --struct UserHandler \
  --package github.com/example/myapp/handler
```

#### ファイルに出力

```bash
convinient_wire analyze \
  --struct UserHandler \
  --package github.com/example/myapp/handler \
  --output dependencies.yaml
```

#### サンプルプロジェクトで実行

```bash
cd sample/basic
convinient_wire analyze \
  --struct UserHandler \
  --package github.com/rmocchy/convinient_wire/sample/basic/handler \
  --output user_handler_deps.yaml
```

## 出力例

```yaml
struct_name: UserHandler
package_path: github.com/rmocchy/convinient_wire/sample/basic/handler
init_functions:
  - name: NewUserHandler
    package_path: github.com/rmocchy/convinient_wire/sample/basic/handler
    signature: func(service github.com/rmocchy/convinient_wire/sample/basic/service.UserService) *github.com/rmocchy/convinient_wire/sample/basic/handler.UserHandler
fields:
  - field_name: service
    type: UserService
    node_type: interface
    package_path: github.com/rmocchy/convinient_wire/sample/basic/service
```

## YAML出力フォーマット

### ルート構造

- `struct_name`: 構造体名
- `package_path`: パッケージパス
- `init_functions`: 初期化関数のリスト
- `fields`: フィールドのリスト
- `skipped`: (オプション) 解析がスキップされた場合true
- `skip_reason`: (オプション) スキップされた理由

### フィールド構造

各フィールドには以下の情報が含まれます:

- `field_name`: フィールド名
- `type`: 型名
- `node_type`: ノードタイプ (`struct`, `interface`, `builtin`)
- `package_path`: パッケージパス (ビルトイン型以外)
- `init_functions`: 初期化関数のリスト (該当する場合)
- `fields`: 再帰的なフィールドのリスト (構造体の場合)

### 初期化関数の構造

- `name`: 関数名
- `package_path`: パッケージパス
- `signature`: 関数シグネチャ

## サンプルプロジェクト

`sample/basic`ディレクトリには、Wire DIパターンを使用したサンプルプロジェクトが含まれています:

```
sample/basic/
├── handler/
│   └── user_handler.go      # UserHandler構造体
├── service/
│   └── user_service.go      # UserServiceインターフェースと実装
├── repository/
│   └── user_repository.go   # UserRepositoryインターフェースと実装
└── wire.go                   # Wire設定
```

## プロジェクト構造

```
.
├── cmd/                      # CLIコマンド (Cobra)
│   ├── root.go
│   └── analyze.go
├── internal/
│   ├── analyze/             # 解析ロジック
│   │   ├── analyze.go
│   │   ├── types.go
│   │   └── yaml_output.go
│   ├── cache/               # キャッシュ実装
│   └── core/                # コア機能
├── sample/                  # サンプルプロジェクト
└── main.go
```

## 内部実装

このツールは以下の技術を使用しています:

- **Cobra**: CLIフレームワーク
- **golang.org/x/tools/go/packages**: Goパッケージの解析
- **gopkg.in/yaml.v3**: YAML出力
- **go/types**: Goの型システムの解析

## ライセンス

[ライセンス情報を追加]

## コントリビューション

プルリクエストやissueは大歓迎です！
