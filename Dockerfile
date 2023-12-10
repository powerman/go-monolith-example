FROM golang:1.16 as prepare

ENV GO111MODULE on
ENV GOBIN "/usr/local/bin"
ENV CGO_ENABLED 0

WORKDIR /app

COPY . .

RUN go install github.com/bufbuild/buf/cmd/buf@v1.1.0
RUN go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

FROM prepare as build

ARG MONO_VERSION="latest"
ENV BUILD_VERSION "${MONO_VERSION}"

WORKDIR /app

RUN go generate ./api/...

RUN go build -ldflags "-X '$(go list -m)/pkg/def.ver=${BUILD_VERSION}'" -o bin/ ./cmd/mono

FROM alpine:3.13 as runner

LABEL org.opencontainers.image.source="https://github.com/powerman/go-monolith-example"

WORKDIR /app

HEALTHCHECK --interval=30s --timeout=5s \
    CMD wget -q -O - http://$HOSTNAME:17000/health-check || exit 1

COPY --from=build "/app/bin/mono" "mono"
COPY --from=build "/app/ms/auth/internal/migrations" "ms/auth/internal/migrations"
COPY --from=build "/app/ms/example/internal/migrations" "ms/example/internal/migrations"

ENTRYPOINT [ "/app/mono" ]

CMD [ "serve" ]
