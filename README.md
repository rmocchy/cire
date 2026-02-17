# Cire

Google Wire の `wire.go` を自動生成する CLI ツール

## インストール

```bash
go install github.com/rmocchy/cire@latest
```

## 使い方

### 1. ルート構造体を定義 (`cire.go`)

```go
package main

import "myapp/handler"

type App struct {
    Handler *handler.UserHandler
}
```

### 2. cire を実行

```bash
cire generate -f ./cire.go
```

### 3. wire を実行

```bash
wire ./
```

## サンプル

- [sample/basic/](sample/basic/)
- [sample/complex/](sample/complex/)
