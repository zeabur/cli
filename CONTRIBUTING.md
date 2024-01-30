# Contributing to `cli`

## Getting Started

See [README.md](./README.md).

## Tests

In `cli`, writing tests is important to ensure the quality and correctness of the code. It can catch bugs early in the development process and prevent regressions when making changes to the codebase.

There are two kinds of tests in `cli`:

- specific pkg test: `cd xxx && go test ./...` or `cd xxx && ginkgo .`
- specific test: `cd xxx && ginkgo -focus xxx` (xxx is `Describe` name)

Once you have written your tests, you should run them before committing your code. For unit tests, you can do this by running the `go test` command in the directory containing your test files, and this command will run all the tests in your package and report any errors. For system tests, you can check if it works by running `./zbpack [folder path]` manually.

## Code Style

- **Write the tests** for every new feature you add.
- Run the tests by running `make test`.
- Format your code by running `gofumpt -w .`. You may need to [install gofumpt](https://github.com/mvdan/gofumpt) before running this command.
- Lint your code before committing by running `make lint`. You may need to [install golangci-lint](https://golangci-lint.run/) before running this command.

## Commit Messages

We use the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format for our commit messages. It makes the commit message clear for readability.

Each commit message should have a type, a scope, and a subject. Here is an example:

```plain
feat(cmd/service): Optimization for user experience

This commit add multiple selectors for user to select services interactively.  
```

The scope can be:

- `cmd`: The command-line interface.
- `api`: The api interface of zeabur cli. (`pkg/api/*`)
- `component`: The library exposed to the users. (`pkg/*`)
- `model`: The model definition of api interface. (`pkg/model/*`)
- `util`: The utility functions. (`internal/utils, pkg/util, intermal/cmdutil`)
- `lint`: The configuration of linters, formatters, `.editerconfig`, etc.
- `docs`: The documentation for zeabur cli. (`docs/*`)
- (feel free to add your own scope if these can not fulfill your changes)

You can contain subscopes in your scope. For example, `cmd/service`.

## Pull Requests

1. Create a new branch for your changes.
2. Make your changes and commit them with clear commit messages following the guidelines above.
3. Push your branch to your fork of the repository.
4. Open a pull request against the `main` branch of the original repository.
5. Wait for a maintainer to review and merge your changes.

Thank you for your contributions!
