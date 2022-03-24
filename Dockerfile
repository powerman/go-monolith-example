FROM alpine:3.15.2

LABEL org.opencontainers.image.source="https://github.com/powerman/go-monolith-example"

WORKDIR /app

HEALTHCHECK --interval=30s --timeout=5s \
    CMD wget -q -O - http://$HOSTNAME:17000/health-check || exit 1

COPY . .

ENTRYPOINT [ "bin/mono" ]

CMD [ "serve" ]
