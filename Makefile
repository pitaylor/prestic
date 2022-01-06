.PHONY: build test clean

PROGRAM = prestic

build: out/$(PROGRAM)-darwin-amd64 out/$(PROGRAM)-linux-amd64

out/$(PROGRAM)-darwin-amd64: $(wildcard *.go)
	GOOS=darwin GOARCH=amd64 go build -o $@ $^

out/$(PROGRAM)-linux-amd64: $(wildcard *.go)
	GOOS=linux GOARCH=amd64 go build -o $@ $^

test:
	go test -v ./...

clean:
	rm -rf out tmp
