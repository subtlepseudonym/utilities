default: build

build:
	mkdir -p bin
	cd cmd && \
	for dir in *; do \
		go build -o "../bin/$$dir" "./$$dir"; \
	done

.PHONY: all test build
