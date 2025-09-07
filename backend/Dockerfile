# builder
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /erebus ./cmd/erebusd


# runtime
FROM alpine:3.18
COPY --from=builder /erebus /usr/local/bin/erebus
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/erebus"]