# Makefile

execute:
	go run .

build:
	go build -o bin/execute.exe

run: build
	./bin/execute.exe
