.PHONY: all build run clean

APP_NAME := build/prerise

all: build

build: clean
	mkdir build && go build -o $(APP_NAME) cmd/main.go

run: build
	./$(APP_NAME)

clean:
	rm -rf build/
	go clean
