# Complex Sample

このサンプルは、より複雑な依存関係のパターンを示しています。

## 特徴

### 1. 複数の@cire構造体
2つの独立したルート構造体があります：
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
