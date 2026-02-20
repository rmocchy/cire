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

.PHONY: test.integrate
test.integrate: ## 統合テストを実行
	make clean.all && make build && make sample.basic && make sample.complex && make sample.duplicate


# サンプルの生成
.PHONY: sample.basic
sample.basic:
	./cire generate -f ./sample/basic/cire.go -j
	wire ./sample/basic

.PHONY: sample.complex
sample.complex: ## 複雑なサンプル（複数ルート・並列依存）のYAMLを生成し、wireでコード生成
	./cire generate -f ./sample/complex/cire.go -j
	wire ./sample/complex

.PHONY: sample.duplicate
sample.duplicate: ## 重複コンストラクタのエラーサンプル（エラーが期待値）
	@if ./cire generate -f ./sample/duplicate/cire.go; then \
		echo "ERROR: Expected failure but succeeded"; \
		exit 1; \
	else \
		echo "OK: Duplicate constructor error detected as expected"; \
	fi

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
	## duplicate
	rm -f ./sample/duplicate/*_di_tree.json

clean.build: ## ビルド成果物をクリーンアップ
	@echo "=== ビルド成果物のクリーンアップ ==="
	rm -f cire
	@echo "✓ ビルド成果物のクリーンアップが完了しました"

# ヘルプ
help: ## このヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
