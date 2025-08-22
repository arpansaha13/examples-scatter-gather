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
    -o ./bin/gateway .

FROM alpine:3.22.0

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/gateway /usr/local/bin/

EXPOSE 4000

CMD ["gateway"]