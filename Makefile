build: mock
	go build -o zc ./cmd/main.go

test:
	go test ./...

mock:
	rm -rf mocks
	mockery --dir=./pkg --output=./mocks/pkg --unroll-variadic=false
