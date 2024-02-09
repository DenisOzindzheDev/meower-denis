FROM golang:latest AS build

WORKDIR /go/src/meower-denis
COPY . .
RUN go mod download
RUN go build -o meow-service cmd/meow-services/meow-service/main.go
RUN go build -o pusher-service cmd/meow-services/pusher-service/main.go
RUN go build -o query-service cmd/meow-services/query-service/main.go

EXPOSE 8080 

FROM alpine:3.7
WORKDIR /usr/bin
COPY --from=build /go/bin .