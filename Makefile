# Server-related targets
build:
	@go build -o babserver ./cmd/server

serve: build
	@-killall babserver 2>/dev/null
	@./babserver&

publish:
	@go run ./cmd/engine

.PHONY: build serve publish
