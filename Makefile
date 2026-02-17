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
sample.basic: ## サンプルファイルの依存関係YAMLを生成
	@echo "=== サンプルファイルの依存関係を生成 ==="
	./cire analyze --file ./sample/basic/cire.go \
	--output ./sample/basic/cire.yaml
	@echo "✓ ./sample/basic/cire.yaml を生成しました"

.PHONY: sample.complex
sample.complex: ## 複雑なサンプル（複数ルート・並列依存）のYAMLを生成
	@echo "=== 複雑なサンプルファイルの依存関係を生成 ==="
	./cire analyze --file ./sample/complex/cire.go \
	--output ./sample/complex/cire.yaml
	@echo "✓ ./sample/complex/cire.yaml を生成しました"

# クリーンアップ
.PHONY: clean
clean: ## ビルド成果物をクリーンアップ
	rm -f convinient_wire
	rm -f coverage.out coverage.html
	rm -f /tmp/user_handler_test.yaml

# ヘルプ
help: ## このヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
