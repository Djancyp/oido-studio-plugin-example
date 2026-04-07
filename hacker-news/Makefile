.PHONY: build clean

PLUGIN_NAME := hacker-news
BINARY := $(PLUGIN_NAME)-plugin

build:
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) .
	@echo "Done! Binary: $(BINARY)"
	@ls -lh $(BINARY)

clean:
	rm -f $(BINARY)

install: build
	@echo "Plugin ready at plugins/$(PLUGIN_NAME)/"
	@echo "Make sure plugin.json is present"
