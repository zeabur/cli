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

### Using Script

> Linux/macOS:

```bash
curl -sSL https://dub.sh/zeabur | sh
```

> Windows

```powershell
iwr -useb https://dub.sh/zb-win | iex
```

**(then restart your powershell)**

### HomeBrew (Linux/macOS)

```bash
brew install zeabur/tap/cli
```

### Scoop (Windows)

```powershell
scoop bucket add zeabur https://github.com/zeabur/cli
scoop install zeabur
```

### 2. Login

If you can open the browser:

```shell
zeabur auth login
```

Or you can use token to login:
```shell
zeabur auth login --token <your-token>
```

Zeabur CLI will open a browser window and ask you to login with your Zeabur account.

### 3. Manage your resources(Interactive mode, recommended)

[![asciicast](https://asciinema.org/a/Olf52EUOCrKU6NGJMbYTw24SL.svg)](https://asciinema.org/a/Olf52EUOCrKU6NGJMbYTw24SL)

```shell
# list all projects
zeabur project ls

# set project context, the following commands will use this project context
# you can use arrow keys to select the project
zeabur context set project

# list all services in the project
zeabur service ls

# set service context(optional)
zeabur context set service

# set environment context(optional)
zeabur context set env

# restart the service
zeabur service restart

# get the latest deployment info
zeabur deployment get

# get the latest deployment log(runtime)
zeabur deployment log -t=runtime

# get the latest deployment log(build)
zeabur deployment log -t=build
```

### 4. Manage your resources(Non-interactive mode)

Non-interactive mode is useful when you want to use Zeabur CLI in a script(such as CI/CD pipeline, etc.)

Note: you can add `-i=false` to all commands to disable interactive mode. 
**In fact, if the parameters are complete, it's same whether you use interactive mode or not.**

```shell
# list all projects
zeabur project ls -i=false

# set project context, the following commands will use this project context
zeabur context set project --name <project-name>
# or you can use project id
# zeabur context set project --id <project-id>

# list all services in the project
zeabur service ls

# set service context(optional)
zeabur context set service --name <service-name>
# or you can use service id
# zeabur context set service --id <service-id>

# set environment context(optional)(only --id is supported)
zeabur context set env --id <env-id>

# restart the service
# if service context is set, you can omit the service name; so does environment context
zeabur service restart --env-id <env-id> --service-name <service-name>
# or you can use service id
# zeabur service restart --env-id <env-id> --service-id <service-id>

# get the latest deployment info(if contexts are set, you can omit the parameters)
zeabur deployment get --env-id <env-id> --service-name <service-name>
# or you can use service id
# zeabur deployment get --env-id <env-id> --service-id <service-id>

# get the latest deployment log(runtime)(service id is also supported)
zeabur deployment log -t=runtime --env-id <env-id> --service-name <service-name>
# get the latest deployment log(build)(service id is also supported)
zeabur deployment log -t=build --env-id <env-id> --service-name <service-name>
```

5. More commands

```shell
zeabur <command> --help
```

## Development Guide

[Development Guide](docs/development_guide.md)

## Acknowledgements

1. GitHub
    * GitHub provides us a place to store the source code of this project and running the CI/CD pipeline.
    * [cli/cli](https://github.com/cli/cli) provides significant inspiration for the organizational structure of this project.
    * [cli/oauth](https://github.com/cli/oauth) we write our own CLI browser OAuth flow based on this project.
