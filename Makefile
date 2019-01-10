.SILENT:
.ONESHELL:
.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.PHONY: run deps build clean exec test trace bench

run: build exec clean

exec:
	./bin/app

build:
	CGO_ENABLED=0 go build -o bin/app -ldflags '-s -w -extldflags "-static"'

test:
	go test -v -count=1 ./pkg/...

trace:
	go tool trace trace

bench:
	go test -benchmem -run=^$ github.com/chneau/solver -bench ^BenchmarkEvaluate$

clean:
	rm -rf bin
	rm -rf upload
	rm -f *.out
	rm -f *.test
	rm -f trace

deps:
	govendor init
	govendor add +e
	govendor update +v

dev:
	go get -u -v github.com/kardianos/govendor

