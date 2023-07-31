# Zeabur CLI

## Overview

Core Features:

1. Login with browser or token
2. Manage your Zeabur resources with CLI
3. The design of the context makes it easier for you to manage services.

## Quick Start

### 1. Install

* Linux/macOS: `curl -sSL https://raw.githubusercontent.com/zeabur/cli/main/hack/install.sh | sh`
* Windows: go to [release page](https://github.com/zeabur/cli/releases) to download the latest version.

(TIP: you can put the binary file in the PATH environment variable to use it conveniently.)

### 2. Login

If you can open the browser:
```shell
./zeabur auth login
```

Or you can use token to login:
```shell
./zeabur auth login --token <your-token>
```

Zeabur CLI will open a browser window and ask you to login with your Zeabur account.

### 3. Manage your resources(Interactive mode, recommended)

```shell
# list all projects
./zeabur project ls

# set project context, the following commands will use this project context
# you can use arrow keys to select the project
./zeabur context set project

# list all services in the project
./zeabur service ls

# set service context(optional)
./zeabur context set service

# set environment context(optional)
./zeabur context set env

# restart the service
./zeabur service restart

# get the latest deployment info
./zeabur deployment get

# get the latest deployment log(runtime)
./zeabur deployment log -t=runtime

# get the latest deployment log(build)
./zeabur deployment log -t=build
```

### 4. Manage your resources(Non-interactive mode)

Non-interactive mode is useful when you want to use Zeabur CLI in a script(such as CI/CD pipeline, etc.)

```shell
# list all projects
./zeabur project ls

# set project context, the following commands will use this project context
./zeabur context set project --name <project-name>
# or you can use project id
# ./zeabur context set project --id <project-id>

# list all services in the project
./zeabur service ls

# set service context(optional)
./zeabur context set service --name <service-name>
# or you can use service id
# ./zeabur context set service --id <service-id>

# set environment context(optional)(only --id is supported)
./zeabur context set env --id <env-id>

# restart the service
# if service context is set, you can omit the service name; so does environment context
./zeabur service restart --env-id <env-id> --name <service-name> 
# or you can use service id
# ./zeabur service restart --env-id <env-id> --id <service-id>

# get the latest deployment info(if contexts are set, you can omit the parameters)
./zeabur deployment get --env-id <env-id> --name <service-name>
# or you can use service id
# ./zeabur deployment get --env-id <env-id> --id <service-id>

# get the latest deployment log(runtime)(service id is also supported)
./zeabur deployment log -t=runtime --env-id <env-id> --name <service-name>
# get the latest deployment log(build)(service id is also supported)
./zeabur deployment log -t=build --env-id <env-id> --name <service-name>
```

5. More commands

```shell
./zeabur help 
```

## Development Guide

[Development Guide](docs/development_guide.md)

## Acknowledgements

1. GitHub
    * GitHub provides us a place to store the source code of this project and running the CI/CD pipeline.
    * [cli/cli](https://github.com/cli/cli) provides significant inspiration for the organizational structure of this project.
    * [cli/oauth](https://github.com/cli/oauth) we write our own CLI browser OAuth flow based on this project.
