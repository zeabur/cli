build: mock
	go build -o zeabur ./cmd/main.go

test:
	go test ./...

mock:
	# if mocks dir exists, remove it
	if [ -d "./mocks" ]; then \
		rm -rf ./mocks; \
	fi
	mockery --dir=./pkg --output=./mocks/pkg --unroll-variadic=false
