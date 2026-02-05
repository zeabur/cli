# Development Guide

## Tools and Concepts You Need to Know

* Go, particularly `interface`.
* [GraphQL](https://graphql.org/) - Query language for the Zeabur API.
* [go-graphql-client](https://github.com/hasura/go-graphql-client) - GraphQL client for Go.
* [Cobra](https://github.com/spf13/cobra) - CLI framework.
* [Viper](https://github.com/spf13/viper) - Configuration framework. We use it to unify environment variables, flags, and configuration files.

Familiarity with the following tools will be beneficial:

* [Ginkgo](https://onsi.github.io/ginkgo/) - Behavior Driven Development (BDD) test framework.

## Development

1. Run from Source Code (Recommended)

Run the command you want, e.g., `go run cmd/main.go auth login --debug`.

2. Build and Run

Use `make build` to build the binary.

Run the binary with `./zeabur auth login --debug`.

## Testing

Run Tests:

* All tests: `make test`.
* Specific package test: `cd xxx && go test ./...` or `cd xxx && ginkgo .`.
* Specific test: `cd xxx && ginkgo -focus xxx` (where xxx is the `Describe` name).

Add Tests (using `internal/cmd/auth/login` as an example):

1. `cd internal/cmd/auth/login`.
2. Run `ginkgo bootstrap` to generate the suite file `login_suite_test.go`.
3. Run `ginkgo generate login` to create the test file `login_test.go`.

## Publishing

1. Tag a new version in the `vx.x.x` format on GitHub.
2. Run `goreleaser build --clean` to build the binaries for Windows, macOS, and Linux.
3. Navigate to the `npm` directory and run `bash prepare.sh` to update the tag and copy the artifacts.
4. Execute `node index.js` to verify that the CLI functions correctly.
5. Use `npm pack` and `npm publish --dry-run` to confirm the package structure is accurate.
6. Run `npm publish` with Zeabur credentials to publish it to NPM.
