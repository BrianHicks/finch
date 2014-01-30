.PHONY: all test deps clean install

all: finch_darwin finch_linux finch_windows

test: deps
	godep go test -v ./...

lint:
	golint *.go

deps:
	go get -v github.com/kr/godep github.com/golang/lint/golint
	godep restore

clean:
	git clean -fx

finch: clean deps
	godep go build -o finch ./cli

install: clean deps finch
	cp finch $(GOPATH)/bin/finch

finch_darwin: clean deps
	GOOS=darwin godep go build -o finch_darwin ./cli

finch_linux: clean deps
	GOOS=linux godep go build -o finch_linux ./cli

finch_windows: clean deps
	GOOS=linux godep go build -o finch_windows ./cli
