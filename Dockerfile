# ------------------------------------------------------------------------------
# Builder Image
# ------------------------------------------------------------------------------
FROM golang:1.14 AS build

WORKDIR /go/src/github.com/figment-networks/oasishub-indexer

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOARCH=amd64
ENV GOOS=linux

RUN \
  GO_VERSION=$(go version | awk {'print $3'}) \
  GIT_COMMIT=$(git rev-parse HEAD) \
  make build
    
# ------------------------------------------------------------------------------
# Target Image
# ------------------------------------------------------------------------------
FROM alpine:3.10 AS release

WORKDIR /app

COPY --from=build /go/src/github.com/figment-networks/oasishub-indexer/oasishub-indexer /app/oasishub-indexer
COPY --from=build /go/src/github.com/figment-networks/oasishub-indexer/migrations /app/migrations

EXPOSE 8081

ENTRYPOINT ["/app/oasishub-indexer"]
