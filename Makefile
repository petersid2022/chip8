build:
	@echo 'Building binary ARCH=amd64 OS=linux'
	CGO_ENABLED=1 CC=gcc GOOS=linux GOARCH=amd64 go build -tags static -ldflags "-s -w" -o chip8 main.go

run: build
	@echo 'Running...'
	./main

clean:
	@echo 'Cleaning...'
	go clean
