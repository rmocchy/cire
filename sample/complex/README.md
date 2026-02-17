# Complex Sample

このサンプルは、より複雑な依存関係のパターンを示しています。

## ファイル構成

- `cire_structs.go`: 依存関係解析対象の構造体定義（`//go:build cire` タグ付き）
- `wire.go`: Google Wireの設定ファイル
- `wire_gen.go`: Wireによって自動生成されるファイル

## 特徴

### 1. 複数のルート構造体
`cire_structs.go` 内に2つの独立したルート構造体があります：
- `UserAppSet`: UserHandlerのみ
- `OrderAppSet`: ProductHandlerとOrderHandler

### 2. 並列依存
`OrderService` は2つのリポジトリに並列で依存しています：
- `UserRepository`
- `ProductRepository`

## 依存関係構造

### UserAppSet (単一ルート)
```
UserAppSet
└── UserHandler → UserService → UserRepository
```

### OrderAppSet (複数ルート + 並列依存)
```
OrderAppSet
├── ProductHandler → ProductService → ProductRepository
└── OrderHandler → OrderService → [UserRepository, ProductRepository] (並列依存)
```

## 実行

```bash
make sample.complex
```

これにより `cire_structs.go` 内のすべての構造体（`UserAppSet` と `OrderAppSet`）の依存関係が `cire.yaml` に出力されます。
