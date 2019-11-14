default: build

build:
	mkdir -p bin
	cd utilities && \
	for dir in *; do \
		go build -o "../bin/$$dir" "./$$dir"; \
	done

.PHONY: all test build
