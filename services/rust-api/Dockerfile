FROM rust:alpine AS builder
WORKDIR /app
RUN apk add --no-cache musl-dev
COPY Cargo.toml ./
RUN mkdir src && echo "fn main() {}" > src/main.rs && cargo build --release && rm -rf src
COPY src ./src
RUN touch src/main.rs && cargo build --release

FROM scratch
COPY --from=builder /app/target/release/rust-api /rust-api
USER 1001:1001
EXPOSE 8081
ENTRYPOINT ["/rust-api"]
