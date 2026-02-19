files := $(wildcard content/posts/*.md)

publish: published.json

published.json: $(files)
	@go run ./cmd/engine -candidates $$(echo $? | tr ' ' ',')

content/posts/%.md::

.PHONY: all clean publish
