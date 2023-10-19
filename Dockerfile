FROM golang:1.21 AS build

COPY . /tmp/flashwise
WORKDIR /tmp/flashwise

RUN go build .

FROM alpine:latest AS runtime

# Fix for "file not found"
RUN apk add --no-cache libc6-compat

COPY --from=build /tmp/flashwise/flashwise /usr/local/bin/flashwise

RUN chmod -R 777 /usr/local/bin/flashwise

ENTRYPOINT [ "/usr/local/bin/flashwise" ]
