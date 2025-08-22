FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

ENV CGO_ENABLED=1
ENV GOCACHE=/go-cache
ENV GOMODCACHE=/gomod-cache

COPY . .

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    go build -ldflags="-w -s" \
    -o ./bin/service_a .

FROM alpine:3.22.0

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/service_a /usr/local/bin/

EXPOSE 4010

CMD ["service_a"]