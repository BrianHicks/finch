.PHONY: all test deps clean

all: finch

test: deps
	godep go test -v ./...

deps:
	go get github.com/kr/godep
	godep restore

clean:
	git clean -fx

finch: clean deps test
	godep go build -o finch ./cli
