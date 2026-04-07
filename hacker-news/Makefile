.PHONY: build clean dist

PLUGIN_NAME := hacker-news
BINARY := $(PLUGIN_NAME)-plugin
DIST_DIR := dist

build:
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) .
	@echo "Done! Binary: $(BINARY)"
	@ls -lh $(BINARY)

# Build and package as a zip ready for upload via the plugin UI
dist: build
	@mkdir -p $(DIST_DIR)
	@echo "Packaging $(PLUGIN_NAME).zip for upload..."
	@cd $(DIST_DIR) && zip -j ../$(PLUGIN_NAME).zip ../plugin.json ../$(BINARY)
	@echo "Done! Upload $(PLUGIN_NAME).zip via the Plugins UI"
	@ls -lh $(PLUGIN_NAME).zip

clean:
	rm -f $(BINARY)
	rm -rf $(DIST_DIR)
	rm -f $(PLUGIN_NAME).zip
