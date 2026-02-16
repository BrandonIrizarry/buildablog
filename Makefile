files := $(wildcard content/posts/*.md)

# Move to /tmp as a kind of "safe delete".
clean:
	mv candidates.txt /tmp

publish: published.json

published.json: $(files)
	@go run ./cmd/engine -candidates $$(echo $? | tr ' ' ',')

content/posts/%.md::

.PHONY: all clean publish
