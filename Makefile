build:
	go build -o zeabur ./cmd/main.go

test:
	go test ./...

mock:
	rm -rf mocks
	mockery --dir=./pkg --output=./mocks/pkg --unroll-variadic=false

lint:
	golangci-lint run ./...
