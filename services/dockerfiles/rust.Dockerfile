FROM rust:1.48.0 as planner
WORKDIR /app
RUN cargo install cargo-chef
COPY . .
RUN cargo chef prepare  --recipe-path recipe.json

FROM rust:1.48.0 as cacher
WORKDIR /app
RUN cargo install cargo-chef
COPY --from=planner /app/recipe.json recipe.json
RUN cargo chef cook --recipe-path recipe.json

FROM rust:1.48.0 as builder
WORKDIR /app
COPY . .
COPY --from=cacher /app/target target
RUN cargo build

FROM debian:buster-slim as runtime
ARG EXE_NAME
RUN apt-get update && apt-get install -y openssl && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/target/debug/${EXE_NAME} /usr/local/bin/myapp
COPY static /static
CMD ["myapp"]
