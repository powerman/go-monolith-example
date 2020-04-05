# Monolith with embedded microservices on Golang


## Setup
1. Install Go 1.14.
2. Install tools required to build/test project:

    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.24.0
    go get github.com/golang/mock/mockgen@v1.4.3


## Test
- `go test ./...` - test project, fast
- `./test -race -tags integration` - thoroughly test project, slow


## Local Environment
TODO В данный момент запускаются только необходимые проекту сторонние
сервисы, сборка и запуск самого проекта будет реализована позднее.

### Requirements
- docker 19.03
- docker-compose 1.25

### Setup
1. После клонирования этого репо скопируйте `env.example.sh` в `env.sh`.
2. Посмотрите `env.sh` в корне этого репо, и, если необходимо, измените
   его для соответствия вашей системе.

### Run
- Всегда загружайте `env.sh` *в каждый терминал*, в котором планируете
  выполнять команды связанные с проектом: `source env.sh`.
    - Если изменился `env.example.sh` то при выполнении `source env.sh` вы
      будете об этом уведомлены, и должны будете вручную перенести
      изменения из `env.example.sh` в `env.sh`, скорректировав их для
      соответствия вашей системе.
    - Рекомендуется прописать в настройках шелла
      `alias dc="if test -f env.sh; then source env.sh; fi && docker-compose"`
      и запускать `dc` вместо `docker-compose`, не беспокоясь больше об `env.sh`.
- Проект запускается и управляется через `docker-compose`.

### Cheatsheet
```sh
dc up -d --remove-orphans               # (пере)запустить все сервисы
dc logs -f -t                           # логи всех сервисов
dc logs -f -t имясервиса                # логи сервиса
dc ps                                   # состояние сервисов
dc restart имясервиса                   # перезапустить сервис
dc exec mysql mysql                     # выполнить команду в контейнере сервиса
dc stop && dc rm -f                     # останавливаем проект
docker volume rm mono_имясервиса        # удаляем данные сервиса
```

Использовать `docker-compose down` не рекомендуется - при этом убивается
сеть докера для этого проекта, и в следующий раз `dc up -d` создаст новую…
через некоторое время доступные докеру сети могут закончиться, и нужно
будет перезапустить сервис докера либо перегрузить комп.
