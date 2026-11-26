# syntax=docker/dockerfile:1

# ---- build stage -----------------------------------------------------------
FROM golang:1.26-alpine AS build

WORKDIR /src

# Cache dependencies before copying the rest of the source.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/api \
    ./cmd/api \
    && CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/migrate \
    ./cmd/migrate

# ---- runtime stage ---------------------------------------------------------
FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app
COPY --from=build /out/api /app/api
COPY --from=build /out/migrate /app/migrate

USER app
EXPOSE 8080

# CMD (instead of ENTRYPOINT) lets `docker compose run --rm backend /app/migrate up`
# replace the default command.
CMD ["/app/api"]
