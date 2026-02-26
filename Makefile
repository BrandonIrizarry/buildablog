# Server-related targets
build:
	@go build -o babserver ./cmd/server

server: build
	@-killall babserver 2>/dev/null
	@./babserver&

publish:
	@go run ./cmd/engine

.PHONY: build server publish
