# Development Guide

## Tools/Concepts you need to know

* Go, especially `interface`.
* [graphql](https://graphql.org/) - Query language for Zeabur API
* [go-graphql-client](github.com/hasura/go-graphql-client) - GraphQL client for Go
* [Cobra](https://github.com/spf13/cobra) - CLI framework
* [Viper](https://github.com/spf13/viper) - Configuration framework. We use it to unify env, flag, config file.

If you also know the following tools, it will be better.

* [GinkGo](https://onsi.github.io/ginkgo/) - "Behavior Driven Development" (BDD) test framework
* [mockery](https://github.com/vektra/mockery) - A mock code autogenerator for Go

## Development

1. Run from source code(Recommended)

`make mock` first to generate mock code.(Regenerate mock code when you change the interface)

then run the cmd you want, e.g. `go run cmd/main.go auth login --debug`

2. Build and run

`make build` to build the binary.

`./zeabur auth login --debug` to run the binary.

## Test

Run Tests:

* all tests: `make test`
* specific pkg test: `cd xxx && go test ./...` or `cd xxx && ginkgo .`
* specific test: `cd xxx && ginkgo -focus xxx` (xxx is `Describe` name)

Add Tests(`ineternal/cmd/auth/login` as example):

1. `cd internal/cmd/auth/login`
2. `ginkgo bootstrap` to generate suite file `login_suite_test.go`
3. `ginkgo generate login` to generate test file `login_test.go`


