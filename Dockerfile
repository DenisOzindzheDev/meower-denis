FROM golang:1.10.2-alpine3.7 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/meower-denis

COPY Gopkg.lock Gopkg.toml ./
COPY vendor vendor
COPY pkg/util pkg/util
COPY internal/stream internal/stream
COPY internal/db internal/db
COPY internal/elastic internal/elastic
COPY internal/models internal/models
COPY cmd/meow-services/meow-service cmd/meow-services/smeow-service
COPY cmd/meow-services/query-service cmd/meow-services/query-service
COPY cmd/meow-services/pusher-service cmd/meow-services/pusher-service

RUN go install ./...

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .