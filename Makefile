build:
	go build -o zeabur ./cmd/main.go

test:
	go test ./...

lint:
	golangci-lint run ./...
