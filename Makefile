

mock:
	rm -rf mocks
	mockery --dir=./pkg --output=./mocks/pkg --unroll-variadic=false
