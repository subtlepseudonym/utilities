default: build

build:
	mkdir -p bin
	cd cmd && \
	for dir in *; do \
		go build -o "../bin/$$dir" "./$$dir"; \
	done

format fmt:
	go fmt -x ./...

test:
	gotest --race -v ./...

.PHONY: all test build
