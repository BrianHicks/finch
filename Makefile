.PHONY: all test deps clean

all: finch

test:
	godep go test -v ./...

lint:
	golint *.go

deps:
	go get -v github.com/kr/godep github.com/golang/lint/golint
	godep restore

clean:
	git clean -fx

finch: clean deps test
	godep go build -o finch ./cli
