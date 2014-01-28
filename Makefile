.PHONY: all test deps clean

all: finch

test: deps
	go test -v ./...

deps:
	go get -d -t -v ./...

clean:
	git clean -fx

finch: clean deps test
	go build -o bin/finch ./cli
