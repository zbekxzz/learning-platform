APP_NAME=platform
MAIN_PATH=./cmd

.PHONY: run build clean tidy migrate-up migrate-down

run:
	go run $(MAIN_PATH)

build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

clean:
	rm -rf bin

tidy:
	go mod tidy

test:
	go test ./...

lint:
	go vet ./...
