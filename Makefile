.PHONY: all
all: build

.PHONY: build
build: hot

hot: $(shell find . -name "*.go")
	go build -o $@ ./cmd

.PHONY: clean

clean:
	-rm hot

.PHONY: test
test:
	go vet ./...
	go test ./...
