# Server-related targets
server:
	@-killall babserver 2>/dev/null
	@go build  -o babserver ./cmd/server
	@./babserver&

publish:
	@go run ./cmd/engine

.PHONY: server publish
