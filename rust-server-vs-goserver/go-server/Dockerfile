FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o server .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
