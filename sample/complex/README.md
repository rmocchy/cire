# Complex Sample

このサンプルは、より複雑な依存関係のパターンを示しています。

## 特徴

### 1. 複数のルートフィールド
`AppSet` 構造体は3つのハンドラーをルートとして持っています：
- `UserHandler`
- `ProductHandler`
- `OrderHandler`

### 2. 並列依存
`OrderService` は2つのリポジトリに並列で依存しています：
- `UserRepository`
- `ProductRepository`

## 依存関係構造

```
AppSet (ルート構造体)
├── UserHandler -> UserService -> UserRepository
├── ProductHandler -> ProductService -> ProductRepository
└── OrderHandler -> OrderService -> [UserRepository, ProductRepository] (並列依存)
```

## 実行

```bash
make sample.complex
```
