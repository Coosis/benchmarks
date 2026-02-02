FROM rust:1.83-alpine AS builder
RUN apk add --no-cache musl-dev
WORKDIR /app
COPY . .
RUN cargo build --release

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/target/release/rust-server .
EXPOSE 8080
CMD ["./rust-server"]
