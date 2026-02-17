# Basic DI Sample

このサンプルは、Wireを使った基本的な依存性注入のパターンを示します。

## 構成

- Repository層: データアクセス
- Service層: ビジネスロジック
- Handler層: HTTPハンドラー

## ファイル構成

- `main.go`: エントリーポイント
- `cire_structs.go`: 依存関係解析対象の構造体定義（`//go:build cire` タグ付き）
- `wire.go`: Wire設定
- `wire_gen.go`: Wireが生成するファイル(自動生成)
- `repository/`: データアクセス層
- `service/`: ビジネスロジック層

## 実行方法

```bash
# 依存関係のインストール
go mod download

# Wireによるコード生成
wire

# 実行
go run .

# 依存関係YAMLの生成
make sample.basic
```

これにより `cire_structs.go` 内のすべての構造体（この場合は `ControllerSet`）の依存関係が `cire.yaml` に出力されます。
- `handler/`: ハンドラー層
