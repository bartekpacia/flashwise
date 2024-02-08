FROM golang:1.22-alpine3.19 AS builder

ARG CGO_ENABLED=0

WORKDIR /tmp/flashwise

# Copy source files necessary to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source files required for build
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go build ./cmd/flashwise

FROM alpine:3.19 AS runtime

COPY --from=builder /tmp/flashwise/flashwise /usr/local/bin/flashwise

RUN chmod -R 777 /usr/local/bin/flashwise

ENTRYPOINT [ "/usr/local/bin/flashwise" ]
