FROM rust:1.78-alpine AS builder

RUN apk add --no-cache musl-dev

WORKDIR /app

COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release && rm -rf src

COPY src ./src
RUN touch src/main.rs
RUN cargo build --release

FROM alpine:3.19

RUN addgroup -g 1000 videoplatform && \
    adduser -u 1000 -G videoplatform -s /bin/sh -D videoplatform

WORKDIR /app

COPY --from=builder /app/target/release/video-platform-rust /app/server

USER videoplatform

EXPOSE 8082

ENTRYPOINT ["./server"]
