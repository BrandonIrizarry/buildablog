files := $(wildcard content/posts/*.md)

# Move to /tmp as a kind of "safe delete".
clean:
	mv candidates.txt /tmp

candidates.txt: $(files)
	@go run ./cmd/engine -candidates $$(echo $? | tr ' ' ',') > $@

content/posts/%.md::

.PHONY: all clean
