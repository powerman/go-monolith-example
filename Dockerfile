FROM alpine:3.12

LABEL org.opencontainers.image.source="https://github.com/powerman/go-monolith-example"

WORKDIR /app

COPY . .

ENTRYPOINT [ "bin/mono" ]

CMD [ "serve" ]
