FROM golang:1.24.4 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main .

FROM debian:bullseye-slim

WORKDIR /usr/local/bin

# Copy binary and assets
COPY --from=builder /build/main .
COPY --from=builder /build/views ./views/
COPY --from=builder /build/public ./public/

RUN chmod +x main

EXPOSE 8080

CMD ["/usr/local/bin/main"]
