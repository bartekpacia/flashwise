FROM golang:1.21 AS build

ARG CGO_ENABLED=0

WORKDIR /tmp/flashwise

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go build ./cmd/flashwise

FROM alpine:3 AS runtime

COPY --from=build /tmp/flashwise/flashwise /usr/local/bin/flashwise

RUN chmod -R 777 /usr/local/bin/flashwise

ENTRYPOINT [ "/usr/local/bin/flashwise" ]
