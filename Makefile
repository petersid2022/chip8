build:
	@echo 'Building binary ARCH=amd64 OS=linux'
	GOARCH=amd64 GOOS=linux go build -o main main.go

run: build
	@echo 'Running...'
	./main

clean:
	@echo 'Cleaning...'
	go clean
