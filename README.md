# Example Go monolith with embedded microservices and The Clean Architecture

[![PkgGoDev](https://pkg.go.dev/badge/github.com/powerman/go-monolith-example)](https://pkg.go.dev/github.com/powerman/go-monolith-example)
[![Go Report Card](https://goreportcard.com/badge/github.com/powerman/go-monolith-example)](https://goreportcard.com/report/github.com/powerman/go-monolith-example)
[![CI/CD](https://github.com/powerman/go-monolith-example/workflows/CI/CD/badge.svg?event=push)](https://github.com/powerman/go-monolith-example/actions?query=workflow%3ACI%2FCD)
[![CircleCI](https://circleci.com/gh/powerman/go-monolith-example.svg?style=svg)](https://circleci.com/gh/powerman/go-monolith-example)
[![Coverage Status](https://coveralls.io/repos/github/powerman/go-monolith-example/badge.svg?branch=master)](https://coveralls.io/github/powerman/go-monolith-example?branch=master)
[![Project Layout](https://img.shields.io/badge/Standard%20Go-Project%20Layout-informational)](https://github.com/golang-standards/project-layout)
[![Release](https://img.shields.io/github/v/release/powerman/go-monolith-example)](https://github.com/powerman/go-monolith-example/releases/latest)

This project shows an example of how to implement monolith with embedded
microservices (a.k.a. modular monolith). This way you'll get many upsides
of monorepo without it complexity and at same time most of upsides of
microservice architecture without some of it complexity.

The embedded microservices use Uncle Bob's "Clean Architecture", check
[Example Go microservice](https://github.com/powerman/go-service-example)
for more details.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Overview](#overview)
  - [Structure of Go packages](#structure-of-go-packages)
  - [Features](#features)
- [Development](#development)
  - [Requirements](#requirements)
  - [Setup](#setup)
  - [Usage](#usage)
    - [Cheatsheet](#cheatsheet)
- [Run](#run)
  - [Docker](#docker)
  - [Source](#source)
- [TODO](#todo)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Overview

### Structure of Go packages

- `api/*` - definitions of own and 3rd-party (in `api/ext-*`)
  APIs/protocols and related auto-generated code
- `cmd/*` - main application(s)
- `internal/*` - packages shared by embedded microservices, e.g.:
  - `internal/config` - configuration (default values, env) shared by
    embedded microservices' subcommands and tests
  - `internal/dom` - domain types shared by microservices (Entities)
- `ms/*` - embedded microservices, with structure:
  - `internal/config` - configuration(s) (default values, env, flags) for
    microservice's subcommands and tests
  - `internal/app` - define interfaces ("ports") for The Clean
    Architecture (or "Ports and Adapters" architecture) and implements
    business-logic
  - `internal/srv/*` - adapters for served APIs/UI
  - `internal/sub` - adapter for incoming events
  - `internal/dal` - adapter for data storage
  - `internal/migrations` - DB migrations (in both SQL and Go)
  - `internal/svc/*` - adapters for accessing external services
- `pkg/*` - helper packages, not related to architecture and
  business-logic (may be later moved to own modules and/or replaced by
  external dependencies), e.g:
  - `pkg/def/` - project-wide defaults
- `*/old/*` - contains legacy code which shouldn't be modified - this code
  is supposed to be extracted from `old/` directories (and refactored to
  follow Clean Architecture) when it'll need any non-trivial modification
  which require testing

### Features

- [X] Project structure (mostly) follow [Standard Go Project Layout](https://github.com/golang-standards/project-layout).
- [X] Strict but convenient golangci-lint configuration.
- [X] Embedded microservices:
  - [X] Well isolated from each other.
  - [X] Can be easily extracted from monolith into separate projects.
  - [X] Share common configuration (both env vars and flags).
  - [X] Each has own CLI subcommands, DB migrations, ports, metrics, …
- [X] Easily testable code (thanks to The Clean Architecture).
- [X] Avoids (and resists to) using global objects (to ensure embedded
  microservices won't conflict on these global objects).
- [X] CLI subcommands support using [cobra](https://github.com/spf13/cobra).
- [X] Graceful shutdown support.
- [X] Configuration defaults can be overwritten by env vars and flags.
- [X] Example JSON-RPC 2.0 API.
- [X] Example tests, both unit and integration.
- [X] Production logging using [structlog](https://github.com/powerman/structlog).
- [X] Production metrics using Prometheus.
- [X] Docker and docker-compose support.
- [X] Smart test coverage report, with optional support for coveralls.io.
- [X] Linters for Dockerfile and shell scripts.
- [X] CI/CD setup for GitHub Actions and CircleCI.

## Development

### Requirements

- Go 1.15
- [Docker](https://docs.docker.com/install/) 19.03+
- [Docker Compose](https://docs.docker.com/compose/install/) 1.25+
- Tools used to build/test project (feel free to install these tools using
  your OS package manager or any other way, but please ensure they've
  required versions; also note these commands will install some non-Go
  tools into `$GOPATH/bin` for the sake of simplicity):

```sh
curl -sSfL https://github.com/hadolint/hadolint/releases/download/v1.18.0/hadolint-$(uname)-x86_64 | install /dev/stdin $(go env GOPATH)/bin/hadolint
curl -sSfL https://github.com/koalaman/shellcheck/releases/download/v0.7.1/shellcheck-v0.7.1.$(uname).x86_64.tar.xz | tar xJf - -C $(go env GOPATH)/bin --strip-components=1 shellcheck-v0.7.1/shellcheck
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.31.0
go get gotest.tools/gotestsum@v0.5.3
go get github.com/golang/mock/mockgen@v1.4.4
go get github.com/cheekybits/genny@master
curl -sSfL https://github.com/go-swagger/go-swagger/releases/download/v0.25.0/swagger_$(uname)_amd64 | install /dev/stdin $(go env GOPATH)/bin/swagger
```

### Setup

1. After cloning the repo copy `env.sh.dist` to `env.sh`.
2. Review `env.sh` and update for your system as needed.
3. It's recommended to add shell alias `alias dc="if test -f env.sh; then
   source env.sh; fi && docker-compose"` and then run `dc` instead of
   `docker-compose` - this way you won't have to run `source env.sh` after
   changing it.

### Usage

To develop this project you'll need only standard tools: `go generate`,
`go test`, `go build`, `docker build`. Provided scripts are for
convenience only.

- Always load `env.sh` *in every terminal* used to run any project-related
  commands (including `go test`): `source env.sh`.
    - When `env.sh.dist` change (e.g. by `git pull`) next run of `source
      env.sh` will fail and remind you to manually update `env.sh` to
      match current `env.sh.dist`.
- `go generate ./...` - do not forget to run after making changes related
  to auto-generated code
- `go test ./...` - test project (excluding integration tests), fast
- `./scripts/test` - thoroughly test project, slow
- `./scripts/test-ci-circle` - run tests locally like CircleCI will do
- `./scripts/cover` - analyse and show coverage
- `./scripts/build` - build docker image and binaries in `bin/`
    - Then use mentioned above `dc` (or `docker-compose`) to run and
      control the project.
    - Access project at host/port(s) defined in `env.sh`.

#### Cheatsheet

```sh
dc up -d --remove-orphans               # (re)start all project's services
dc logs -f -t                           # view logs of all services
dc logs -f SERVICENAME                  # view logs of some service
dc ps                                   # status of all services
dc restart SERVICENAME
dc exec SERVICENAME COMMAND             # run command in given container
dc stop && dc rm -f                     # stop the project
docker volume rm PROJECT_SERVICENAME    # remove some service's data
```

It's recommended to avoid `docker-compose down` - this command will also
remove docker's network for the project, and next `dc up -d` will create a
new network… repeat this many enough times and docker will exhaust
available networks, then you'll have to restart docker service or reboot.

## Run

### Docker

```
$ docker run -i -t --rm ghcr.io/powerman/go-monolith-example -v
```

### Source

Use of the `./scripts/build` script is optional (it's main feature is
embedding git version into compiled binary), you can use usual
`go get|install|build` to get the application instead.

```
$ ./scripts/build
$ ./bin/mono -h
$ ./bin/mono serve -h
$ ./bin/mono -v
$ ./bin/mono serve
```

## TODO

- [ ] Add gRPC service example.
- [ ] Add OpenAPI service example.
- [ ] Add NATS/STAN publish/subscribe example.
- [ ] Add DAL implementation for Postgresql.
- [ ] Add LPC (local procedure call API between embedded microservices).
- [ ] Embed https://github.com/powerman/go-service-example as an example
  of embedding microservices from another repo.
