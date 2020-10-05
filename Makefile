default: build

build:
	mkdir -p bin
	cd cmd && \
	for dir in *; do \
		go build -o "../bin/$$dir" "./$$dir"; \
	done

format fmt:
	gofmt -l -w .

test:
	gotest --race -v ./...

.PHONY: build format fmt test
