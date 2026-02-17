# デフォルトターゲット
.DEFAULT_GOAL := help

# ビルド
.PHONY: build
build: ## バイナリをビルド
	go build -o cire

# 単体テスト
.PHONY: test
test: ## Goの単体テストを実行
	go test -v ./...

# テストカバレッジ
.PHONY: test.coverage
test.coverage: ## テストカバレッジを表示
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "カバレッジレポートを coverage.html に生成しました"

# サンプルの生成
.PHONY: sample.basic
sample.basic:
	./cire analyze --file ./sample/basic/cire.go
	wire ./sample/basic

.PHONY: sample.complex
sample.complex: ## 複雑なサンプル（複数ルート・並列依存）のYAMLを生成し、wireでコード生成
	./cire analyze --file ./sample/complex/cire.go
	wire ./sample/complex

# クリーンアップ
.PHONY: clean.all clean.sample clean.build
clean.all: ## すべてのビルド成果物をクリーンアップ
	make clean.sample
	make clean.build

clean.sample: ## ビルド成果物をクリーンアップ
	## basic
	rm -f ./sample/basic/cire.yaml
	rm -f ./sample/basic/wire.go
	rm -f ./sample/basic/wire_gen.go
	## complex
	rm -f ./sample/complex/cire.yaml
	rm -f ./sample/complex/wire.go
	rm -f ./sample/complex/wire_gen.go

clean.build: ## ビルド成果物をクリーンアップ
	@echo "=== ビルド成果物のクリーンアップ ==="
	rm -f cire
	@echo "✓ ビルド成果物のクリーンアップが完了しました"

# ヘルプ
help: ## このヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
