.PHONY: clean
clean:
	rm -rf bin
	go clean ./...

.PHONY: setup
setup:
	go install github.com/jessevdk/go-assets-builder

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build: generate
	go build -o bin/server .

.PHONY: test
test:
	go test -v -p 1 ./...
