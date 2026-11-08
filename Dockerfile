
# 1: Build
FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/server ./cmd/server

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/seeder ./cmd/seeder


#2: Runtime
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/bin/seeder /app/seeder
COPY --from=builder /app/db/seeds /app/db/seeds

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 50051

ENTRYPOINT ["/app/server"]

