.PHONY: all
all: build

.PHONY: build
build: hot

hot: $(shell find . -name "*.go") go.mod go.sum
	go build -mod=vendor -o $@ ./cmd

.PHONY: clean

clean:
	-rm hot

.PHONY: test
test:
	go vet ./...
	go test ./...
