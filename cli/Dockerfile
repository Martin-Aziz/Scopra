FROM rust:1.82-alpine AS build
RUN apk add --no-cache musl-dev
WORKDIR /app
COPY Cargo.toml ./
COPY src ./src
RUN cargo build --release

FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app
COPY --from=build /app/target/release/nexus-cli /usr/local/bin/nexus
USER app
ENTRYPOINT ["/usr/local/bin/nexus"]
