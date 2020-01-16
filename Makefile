all: run

run: install
	chuck

fmt:
	go fmt ./...

test: fmt
	go test ./... -count=1

install: fmt
	go install ./...