# Server-related targets
build:
	@go build -o babserver ./cmd/server

serve: build
	@-killall babserver 2>/dev/null
	@./babserver&

.PHONY: build serve
