files := $(wildcard content/posts/*.md)

publish: published.json

published.json: $(files)
	@go run ./cmd/engine -candidates $$(echo $? | tr ' ' ',')

content/posts/%.md::

# Server-related targets
server:
	@-killall babserver 2>/dev/null
	@go build  -o babserver ./cmd/server
	@./babserver&

.PHONY: all clean publish server
