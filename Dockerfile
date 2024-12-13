FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o /app/holders .


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/holders /app/holders
ENTRYPOINT ["/app/holders"]