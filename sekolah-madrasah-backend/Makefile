# AI CDN Services Makefile

.PHONY: swagger swagger-local swagger-undo update-version help

help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

.PHONY: swagger update-version help
swagger: ## Generate swagger documentation (auto-installs swag if needed)
	@echo "Generating swagger documentation..."
	@if command -v swag >/dev/null 2>&1; then \
		swag init; \
	else \
		echo "swag not found. Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		if command -v swag >/dev/null 2>&1; then \
			swag init; \
		elif [ -f ~/go/bin/swag ]; then \
			~/go/bin/swag init; \
		else \
			echo "Failed to install or find swag binary"; \
			exit 1; \
		fi; \
	fi

update-version: ## Update version timestamp in health endpoint with WIB timezone
	@echo "Updating version timestamp in health endpoint..."
	@VERSION_WIB=$$(TZ='Asia/Jakarta' date '+%d %B %Y - %H:%M:%S WIB'); \
	sed -i '' -e 's/Version: *"[^"]*"/Version: "'"$$VERSION_WIB"'"/' config/init_helpers.go
	@echo "✅ Version updated to: $$(TZ='Asia/Jakarta' date '+%d %B %Y - %H:%M:%S WIB')"

swagger-local: ## Generate swagger documentation for local development
	@echo "Generating swagger documentation for local development..."
	@cp main.go main.go.bak
	@REST_PORT=$$(grep '^REST_PORT=' .env | cut -d'=' -f2); \
	sed -i '' "s|// @host .*|// @host localhost:$$REST_PORT|" main.go; \
	sed -i '' "s|// @BasePath .*|// @BasePath /|" main.go; \
	sed -i '' "s|// @schemes .*|// @schemes http|" main.go; \
	if command -v swag >/dev/null 2>&1; then \
		swag init; \
	else \
		echo "swag not found. Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		if command -v swag >/dev/null 2>&1; then \
			swag init; \
		elif [ -f ~/go/bin/swag ]; then \
			~/go/bin/swag init; \
		else \
			echo "Failed to install or find swag binary"; \
			exit 1; \
		fi; \
	fi; \
	mv main.go.bak main.go
	@echo "✅ Swagger documentation generated for local development"

swagger-undo: ## Undo local changes in docs folder (revert to git state)
	@echo "Undoing local changes in docs folder..."
	@if [ -d "docs" ]; then \
		git checkout -- docs/ 2>/dev/null && echo "✅ Successfully reverted docs folder to git state" || \
		echo "ℹ️  No changes to revert in docs folder or not a git repository"; \
	else \
		echo "⚠️  docs folder not found"; \
	fi