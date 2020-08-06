.PHONY: clean
clean:
	rm -rf bin
	go clean ./...

.PHONY: build
build:
	go build -o bin/server .

.PHONY: test
test:
	go test -v -p 1 ./...
