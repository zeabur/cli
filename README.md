# Zeabur CLI

[Zeabur](https://zeabur.com/)'s official command line tool

> Note: Zeabur CLI is currently in beta, and we are still working on it. If you have any questions or suggestions, please feel free to contact us.

## How cool it is

1. Manage your Zeabur resources with CLI
2. Login with browser or token
3. Intuitive and easy to use
4. The design of the context makes it easier for you to manage services.
5. The seamless integration of interactive and non-interactive modes.

## Quick Start

### 1. Install

No need to install, you can use it directly with npx. Make sure you have Node.js installed.

### 2. Login

If you can open the browser:

```shell
npx zeabur auth login
```

Or you can use token to login:
```shell
npx zeabur auth login --token <your-token>
```

Zeabur CLI will open a browser window and ask you to login with your Zeabur account.

### 3. Manage your resources(Interactive mode, recommended)

[![asciicast](https://asciinema.org/a/Olf52EUOCrKU6NGJMbYTw24SL.svg)](https://asciinema.org/a/Olf52EUOCrKU6NGJMbYTw24SL)

```shell
# list all projects
npx zeabur project ls

# set project context, the following commands will use this project context
# you can use arrow keys to select the project
npx zeabur context set project

# list all services in the project
npx zeabur service ls

# set service context(optional)
npx zeabur context set service

# set environment context(optional)
npx zeabur context set env

# restart the service
npx zeabur service restart

# get the latest deployment info
npx zeabur deployment get

# get the latest deployment log(runtime)
npx zeabur deployment log -t=runtime

# get the latest deployment log(build)
npx zeabur deployment log -t=build
```

### 4. Manage your resources(Non-interactive mode)

Non-interactive mode is useful when you want to use Zeabur CLI in a script(such as CI/CD pipeline, etc.)

Note: you can add `-i=false` to all commands to disable interactive mode. 
**In fact, if the parameters are complete, it's same whether you use interactive mode or not.**

```shell
# list all projects
npx zeabur project ls -i=false

# set project context, the following commands will use this project context
npx zeabur context set project --name <project-name>
# or you can use project id
# npx zeabur context set project --id <project-id>

# list all services in the project
npx zeabur service ls

# set service context(optional)
npx zeabur context set service --name <service-name>
# or you can use service id
# npx zeabur context set service --id <service-id>

# set environment context(optional)(only --id is supported)
npx zeabur context set env --id <env-id>

# restart the service
# if service context is set, you can omit the service name; so does environment context
npx zeabur service restart --env-id <env-id> --service-name <service-name>
# or you can use service id
# npx zeabur service restart --env-id <env-id> --service-id <service-id>

# get the latest deployment info(if contexts are set, you can omit the parameters)
npx zeabur deployment get --env-id <env-id> --service-name <service-name>
# or you can use service id
# npx zeabur deployment get --env-id <env-id> --service-id <service-id>

# get the latest deployment log(runtime)(service id is also supported)
npx zeabur deployment log -t=runtime --env-id <env-id> --service-name <service-name>
# get the latest deployment log(build)(service id is also supported)
npx zeabur deployment log -t=build --env-id <env-id> --service-name <service-name>
```

5. More commands

```shell
npx zeabur <command> --help
```

## Workspaces (personal / team)

By default, the CLI acts under the personal workspace — the account that logged in. To list or create projects under a team you belong to, switch the workspace:

```shell
# show your personal workspace + all teams you belong to, with your role per team
npx zeabur workspace list

# switch to a team — pass either the team name OR its 24-char ObjectID.
# A 24-char hex value is always interpreted as an ID; anything else is
# looked up by name. If multiple teams share a name the CLI errors out and
# prints the per-candidate `workspace switch <id>` invocation so you can
# pick by ID (team names are unconstrained, so duplicates are possible).
npx zeabur workspace switch acme

# show the workspace the CLI is currently using
npx zeabur workspace current

# return to the personal workspace
npx zeabur workspace clear
```

Switching a workspace clears the pinned project / environment / service context, because resource IDs do not overlap between workspaces.

The workspace only affects directory-level commands (`project list`, `project create`, `deploy` without a linked project). Commands that take a specific service or deployment ID use that resource's own owner and are workspace-independent — your team's `service restart` works the same regardless of which workspace is active.

For one-off commands that should run under a different workspace without switching the persisted state, use the `--workspace` flag:

```shell
# list projects in the "acme" team without switching workspaces
npx zeabur --workspace acme project list
```

`switch personal` is **not** a way to return to personal — it always looks for a team literally named `personal` (team names are unconstrained). Use `workspace clear` to go back.

## Development Guide

[Development Guide](docs/development_guide.md)

## Acknowledgements

1. GitHub
    * GitHub provides us a place to store the source code of this project and running the CI/CD pipeline.
    * [cli/cli](https://github.com/cli/cli) provides significant inspiration for the organizational structure of this project.
    * [cli/oauth](https://github.com/cli/oauth) we write our own CLI browser OAuth flow based on this project.
