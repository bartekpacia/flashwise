FROM golang:1.21 AS build

WORKDIR /tmp/flashwise

COPY go.mod go.sum ./

RUN go mod download

COPY *.go /tmp/flashwise

RUN go build .

FROM alpine:latest AS runtime

# Fix for "file not found"
RUN apk add --no-cache libc6-compat

COPY --from=build /tmp/flashwise/flashwise /usr/local/bin/flashwise

RUN chmod -R 777 /usr/local/bin/flashwise

ENTRYPOINT [ "/usr/local/bin/flashwise" ]
