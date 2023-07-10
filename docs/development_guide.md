# Development Guide

## Completed

* [x] project structure
* [x] auth workflow
* [x] graphql example

## Need more consideration/discussion

### 1. interactive design

If you want to restart a service, user should tell the CLI which service to restart.
It is not a good idea to ask user to input the service name/id, and environment name/id.

We should provide a interactive way to let user select the service and environment.
In the other word, this feature sounds like "restart current service in current environment".

There is no example in this project now, we need to discuss and design it.

(And there are many other interactive features, e.g. select a project, select a service, select a environment, select a deployment, etc.)

### 2. GraphQL model import recursively

Now we define the GraphQL model in `pkg/graphql/model.go`, models are imported recursively.

Maybe we shouldn't re-use them to avoid the dependency between models, or graphql will query them recursively.

I think it is better to re-define the different models in everywhere we need.

## ToDos

* Business Features
  * [ ] auth
  * [ ] project
  * [ ] service
  * [ ] environment
  * [ ] deployment
  * [ ] billing
  * [ ] templates
* [ ] External Service
  * Blocked by Zeabur Backend
* [ ] Prompt implementation
* Documentation
  * [ ] README.md
  * [ ] docs/development_guide.md
* DevOps
  * [ ] Makefile
  * [ ] CI/CD
  * [ ] Release(go-releaser)
  * [ ] Download Script(or Homebrew, Scoop, etc.)
* [ ] Pretty Print
* golangci-lint
  * [ ] Lint itself(lint config, CI/CD)
  * [ ] Lint the code
* Tests
  * [ ] Not very important but very difficult, welcome our heroes!

And, all `TODO` in the code.

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

`./zc auth login --debug` to run the binary.

## Test

Run Tests:

* all tests: `make test`
* specific pkg test: `cd xxx && go test ./...` or `cd xxx && ginkgo .`
* specific test: `cd xxx && ginkgo -focus xxx` (xxx is `Describe` name)

Add Tests(`ineternal/cmd/auth/login` as example):

1. `cd internal/cmd/auth/login`
2. `ginkgo bootstrap` to generate suite file `login_suite_test.go`
3. `ginkgo generate login` to generate test file `login_test.go`


