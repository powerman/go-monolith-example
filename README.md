# Monolith with embedded microservices on Golang (example)


## Setup
1. Install Go 1.14.
2. Install tools required to build/test project:

```
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.24.0
go get github.com/golang/mock/mockgen@v1.4.3
go get github.com/cheekybits/genny@master
```

## Test
- `go test ./...` - test project, fast
- `./test -race -tags integration` - thoroughly test project, slow


## Local Environment
TODO At the moment docker-compose runs only external dependencies,
building and running monolith itself will be implemented a bit later.

### Requirements
- docker 19.03
- docker-compose 1.25

### Setup
1. After cloning the repo copy `env.example.sh` to `env.sh`.
2. Review `env.sh` and update for your system as needed.

### Run
- Always load `env.sh` *in every terminal* used to run project-related
  commands: `source env.sh`.
    - When `env.example.sh` change (e.g. by `git pull`) next run of
      `source env.sh` will fail and remind you to manually update `env.sh`
      to match current `env.example.sh`.
    - It's recommended to add shell alias `alias dc="if test -f env.sh;
      then source env.sh; fi && docker-compose"` and then run `dc` instead
      of `docker-compose` - this way you won't have to bother about
      `source env.sh` anymore.
- Use `docker-compose` to run and control the project.

### Cheatsheet
```sh
dc up -d --remove-orphans               # (re)start all project's services
dc logs -f -t                           # view logs of all services
dc logs -f -t SERVICENAME               # view logs of some service
dc ps                                   # status of all services
dc restart SERVICENAME
dc exec mysql mysql                     # run command in given container
dc stop && dc rm -f                     # stop the project
docker volume rm mono_SERVICENAME       # remove some service's data
```

It's recommended to avoid `docker-compose down` - this command will also
remove docker's network for the project, and next `dc up -d` will create a
new network… repeat this many enough times and docker will exhaust
available networks, then you'll have to restart docker service or reboot.
