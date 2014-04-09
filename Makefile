.PHONY:  test deps install xc relase

test: deps
	godep go test -parallel=8 -v ./...

lint:
	golint *.go

deps:
	go get -v github.com/kr/godep github.com/golang/lint/golint
	godep restore

xc:
	go get -v github.com/laher/goxc
	goxc -d $(shell pwd)/download-page -pv=$(shell grep -oe '\d\+\.\d\+\.\d\+' main.go | head -n 1) validate compile package

release: xc
	goxc -d $(shell pwd)/download-page -pv=$(shell grep -oe '\d\+\.\d\+\.\d\+' main.go | head -n 1) bintray
