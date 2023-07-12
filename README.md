# Zeabur CLI

## Overview

Core Features:

1. Login with browser or token
2. Manage your Zeabur resources with CLI
3. The design of the context makes it easier for you to manage services.

## Usage

1. login: `./zeabur auth login`
2. manage your resources, such as `./zeabur project ls`, `./zeabur service get`

Tips: you could use `./zeabur context set <context-type>` to set the context.

1. set project context: `./zeabur context set project`
2. set environment context: `./zeabur cpntext set env`
3. set service context: `./zeabur context set service`

## Development Guide

[Development Guide](docs/development_guide.md)

## Acknowledgements

1. GitHub
    * GitHub provides us a place to store the source code of this project and running the CI/CD pipeline.
    * [cli/cli](https://github.com/cli/cli) provides significant inspiration for the organizational structure of this project.
    * [cli/oauth](https://github.com/cli/oauth) we write our own CLI browser OAuth flow based on this project.
